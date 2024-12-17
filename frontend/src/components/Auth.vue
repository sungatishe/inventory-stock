<template>
  <div class="auth-container">
    <div class="form-card">
      <h2>{{ isLogin ? 'Login' : 'Register' }}</h2>
      <form @submit.prevent="handleSubmit">
        <div class="form-group">
          <label for="email">Email</label>
          <input
              type="email"
              id="email"
              v-model="form.email"
              required
              placeholder="Enter your email"
          />
        </div>

        <div class="form-group">
          <label for="password">Password</label>
          <input
              type="password"
              id="password"
              v-model="form.password"
              required
              placeholder="Enter your password"
          />
        </div>

        <div v-if="!isLogin" class="form-group">
          <label for="confirm-password">Confirm Password</label>
          <input
              type="password"
              id="confirm-password"
              v-model="form.confirmPassword"
              required
              placeholder="Confirm your password"
          />
        </div>

        <button type="submit" class="btn-primary">
          {{ isLogin ? 'Login' : 'Register' }}
        </button>
      </form>

      <p class="toggle-link">
        {{ isLogin ? "Don't have an account?" : "Already have an account?" }}
        <a @click.prevent="toggleForm">{{ isLogin ? 'Register' : 'Login' }}</a>
      </p>

      <p class="status-message" v-if="statusMessage">{{ statusMessage }}</p>
    </div>
  </div>
</template>

<script>
import axios from "axios";
import Cookies from "js-cookie";



export default {
  data() {
    return {
      isLogin: true,
      form: {
        email: "",
        password: "",
        confirmPassword: "",
      },
      statusMessage: "",
    };
  },
  methods: {
    toggleForm() {
      this.isLogin = !this.isLogin;
      this.statusMessage = "";
      this.form.email = "";
      this.form.password = "";
      this.form.confirmPassword = "";
    },
    async handleSubmit() {
      const { email, password, confirmPassword } = this.form;
      if (!this.isLogin && password !== confirmPassword) {
        this.statusMessage = "Passwords do not match.";
        return;
      }

      try {
        const url = `http://127.0.0.1:8080/${this.isLogin ? 'login' : 'register'}`;
        const payload = { email, password };

        const response = await axios.post(url, payload);

        if (response.status === 200) {
          this.statusMessage = this.isLogin
              ? "Login successful!"
              : "Registration successful! You can now login.";

          if (this.isLogin && response.data.token) {
            // Сохраняем токен в куки
            Cookies.set('auth_token', response.data.token, { expires: 7, secure: true });
            console.log("Token saved to cookies:", response.data.token);

            this.router.push('/items');
          }
        } else {
          this.statusMessage = response.data.message || "An error occurred.";
        }
      } catch (error) {
        this.statusMessage = error.response?.data?.message || "Server error.";
      }
    },
  },
};
</script>

<style scoped>
.auth-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: #f3f4f6;
  font-family: Arial, sans-serif;
}

.form-card {
  background: white;
  padding: 2rem;
  border-radius: 10px;
  box-shadow: 0 4px 10px rgba(0, 0, 0, 0.1);
  width: 400px;
  text-align: center;
}

h2 {
  margin-bottom: 1rem;
}

.form-group {
  margin-bottom: 1rem;
  text-align: left;
}

label {
  display: block;
  margin-bottom: 0.5rem;
  font-size: 0.9rem;
  color: #333;
}

input {
  width: 100%;
  padding: 0.7rem;
  font-size: 1rem;
  border: 1px solid #ddd;
  border-radius: 5px;
  box-sizing: border-box;
}

input:focus {
  border-color: #3498db;
  outline: none;
}

.btn-primary {
  width: 100%;
  padding: 0.7rem;
  font-size: 1rem;
  background: #3498db;
  color: white;
  border: none;
  border-radius: 5px;
  cursor: pointer;
}

.btn-primary:hover {
  background: #2980b9;
}

.toggle-link {
  margin-top: 1rem;
}

a {
  color: #3498db;
  cursor: pointer;
}

a:hover {
  text-decoration: underline;
}

.status-message {
  margin-top: 1rem;
  color: red;
}
</style>

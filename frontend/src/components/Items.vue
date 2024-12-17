<template>
  <div class="items-container">
    <h2>Your Inventory</h2>

    <!-- Добавление нового элемента -->
    <div class="add-item-form">
      <input
          v-model="newItem.name"
          placeholder="Item name"
          required
      />
      <input
          v-model="newItem.quantity"
          type="number"
          placeholder="Quantity"
          required
      />
      <input
          v-model="newItem.price"
          type="number"
          step="0.01"
          placeholder="Price"
          required
      />
      <button @click="addItem" class="btn-primary">Add Item</button>
    </div>

    <!-- Список элементов -->
    <div class="items-list">
      <div
          v-for="item in items"
          :key="item.id"
          class="item-card"
      >
        <h3>{{ item.name }}</h3>
        <p>Quantity: {{ item.quantity }}</p>
        <p>Price: ${{ item.price }}</p>
        <div class="item-actions">
          <button class="btn-edit" @click="editItem(item)">Edit</button>
          <button class="btn-delete" @click="deleteItem(item.id)">Delete</button>
        </div>
      </div>
    </div>

    <!-- Модальное окно редактирования -->
    <div v-if="editMode" class="modal-overlay">
      <div class="modal-content">
        <h3>Edit Item</h3>
        <input v-model="editItemData.name" placeholder="Item name" />
        <input v-model="editItemData.quantity" type="number" placeholder="Quantity" />
        <input v-model="editItemData.price" type="number" step="0.01" placeholder="Price" />
        <button class="btn-primary" @click="updateItem">Save</button>
        <button class="btn-cancel" @click="editMode = false">Cancel</button>
      </div>
    </div>
  </div>
</template>

<script>
import axios from "axios";
import Cookies from "js-cookie";
import { jwtDecode } from "jwt-decode";

export default {
  data() {
    return {
      items: [], // Список инвентаря
      newItem: { name: "", quantity: 0, price: 0 },
      editMode: false,
      editItemData: { id: null, name: "", quantity: 0, price: 0 },
      userID: "", // User ID из токена
    };
  },
  methods: {
    // Извлекает user_id из токена и сразу загружает список инвентаря
    extractUserIDAndFetchItems() {
      const token = Cookies.get("auth_token");
      if (token) {
        try {
          const decoded = jwtDecode(token);
          this.userID = decoded.user_id?.toString() || "";

          if (this.userID) {
            this.fetchItems();
          } else {
            console.error("User ID missing in token");
          }
        } catch (err) {
          console.error("Failed to decode token:", err);
        }
      } else {
        console.error("No auth token found.");
      }
    },

    // Получаем список элементов из API
    async fetchItems() {
      const token = Cookies.get("auth_token");
      if (!this.userID) return; // Prevent fetching if userID is not available

      try {
        const response = await axios.get(
            `http://127.0.0.1:8080/inventory/${this.userID}`,
            { headers: { Authorization: `Bearer ${token}` } }
        );

        console.log("API Response:", response.data);

        // Filter items by checking for required fields
        this.items = response.data.items.filter(
            item => item.name && item.quantity !== undefined && item.price !== undefined
        );

      } catch (error) {
        console.error("Error fetching items:", error);
      }
    },

    // Добавляем новый элемент
    async addItem() {
      const token = Cookies.get("auth_token");
      if (!this.userID) return;

      const payload = {
        item: {
          ...this.newItem,
          user_id: this.userID, // ID пользователя из токена
        },
      };

      try {
        await axios.post("http://127.0.0.1:8080/inventory", payload, {
          headers: { Authorization: `Bearer ${token}` },
        });
        this.fetchItems(); // Перезагружаем список
        this.newItem = { name: "", quantity: 0, price: 0 };
      } catch (error) {
        console.error("Error adding item:", error);
      }
    },

    // Удаляем элемент
    async deleteItem(itemID) {
      const token = Cookies.get("auth_token");
      if (!this.userID) return;

      try {
        await axios.delete(
            `http://127.0.0.1:8080/inventory/${this.userID}/${itemID}`,
            { headers: { Authorization: `Bearer ${token}` } }
        );
        this.fetchItems();
      } catch (error) {
        console.error("Error deleting item:", error);
      }
    },

    // Подготавливаем элемент для редактирования
    editItem(item) {
      this.editItemData = { ...item };
      this.editMode = true;
    },

    // Сохраняем изменения элемента
    async updateItem() {
      const token = Cookies.get("auth_token");
      if (!this.userID) return;

      const payload = {
        item: {
          ...this.editItemData,
          user_id: this.userID,
        },
      };

      try {
        await axios.put("http://127.0.0.1:8080/inventory", payload, {
          headers: { Authorization: `Bearer ${token}` },
        });
        this.fetchItems();
        this.editMode = false;
      } catch (error) {
        console.error("Error updating item:", error);
      }
    },
  },
  mounted() {
    this.extractUserIDAndFetchItems();
  },
};
</script>

<style scoped>
.items-container {
  padding: 2rem;
  text-align: center;
}

.add-item-form {
  margin-bottom: 1.5rem;
}

.add-item-form input {
  margin: 0.5rem;
  padding: 0.5rem;
  border: 1px solid #ddd;
  border-radius: 5px;
}

.items-list {
  display: flex;
  flex-wrap: wrap;
  gap: 1rem;
  justify-content: center;
}

.item-card {
  background-color: #f9f9f9;
  border: 1px solid #ddd;
  border-radius: 5px;
  padding: 1rem;
  width: 250px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.item-actions {
  display: flex;
  justify-content: space-between;
  margin-top: 1rem;
}

.btn-primary,
.btn-edit,
.btn-delete {
  padding: 0.5rem 1rem;
  border: none;
  color: white;
  cursor: pointer;
  border-radius: 5px;
}

.btn-primary {
  background-color: #3498db;
}

.btn-edit {
  background-color: #f1c40f;
}

.btn-delete {
  background-color: #e74c3c;
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
}

.modal-content {
  background: white;
  padding: 2rem;
  border-radius: 10px;
  width: 300px;
}

.btn-cancel {
  margin-left: 0.5rem;
  background: gray;
  color: white;
}
</style>

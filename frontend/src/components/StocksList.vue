<template>
  <div class="p-4">
    <h1 class="text-3xl font-bold mb-4 text-green-600 text-center title-font">MyStocks</h1>
    <div v-if="loading" class="text-center">Loading...</div>
    <div v-else-if="error" class="text-red-500">{{ error }}</div>

    <table v-else class="table-auto w-full border-collapse max-w-[800px] mx-auto">
      <thead>
      <tr>
        <th class="p-2 border-b border-gray-300 text-left">Ticker</th>
        <th class="p-2 border-b border-gray-300 text-left">Company Name</th>
        <th class="p-2 border-b border-gray-300 text-right">Current Price</th>
        <th class="p-2 border-b border-gray-300">Recommendation</th>
      </tr>
      </thead>
      <tbody>
      <StockRow
          v-for="stock in stocks"
          :key="stock.ticker"
          :stock="stock"
          @select="selectStock"
      />
      </tbody>
    </table>

    <StockDetailsModal
        v-if="selectedStock"
        :stock="selectedStock"
        @close="selectedStock = null"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import axios from 'axios';
import StockDetailsModal from "@/components/StockDetailsModal.vue";
import StockRow from "@/components/StockRow.vue";

const stocks = ref([]);
const loading = ref(true);
const error = ref<string | null>(null);
const selectedStock = ref<any>(null);

const apiBaseUrl = import.meta.env.VITE_STOCKS_API_URL;

onMounted(async () => {
  try {
    console.log(`${apiBaseUrl}/stocks`);
    const response = await axios.get(`${apiBaseUrl}/stocks`);
    console.log(response.data);
    stocks.value = response.data;
  } catch (e) {
    error.value = 'Failed to load stocks';
  } finally {
    loading.value = false;
  }
});

const selectStock = async (ticker: string) => {
  try {
    const response = await axios.get(`${apiBaseUrl}/stocks/${ticker}`);
    selectedStock.value = response.data;
  } catch (e) {
    error.value = 'Failed to load stock details';
  }
};

</script>

<style>
body {
  font-family: 'Inter', sans-serif;
  background-color: #f9fafb
}
/* Title styling */
h1.title-font {
  font-family: 'Poppins', sans-serif;
  letter-spacing: 1px;
  font-size: clamp(1.5rem, 4vw, 2.5rem);
}

/* Table styling */
table {
  font-family: 'Inter', sans-serif;
  color: #4b5563;
  width: 100%;
  border-collapse: collapse;
}

/* Table header styling */
th {
  font-weight: 600;
  color: #374151;
  padding: 0.5rem;
}

/* Table rows hover effect */
tr:hover {
  background-color: #f3f4f6;
  transition: background-color 0.2s ease-in-out;
}

/* Responsive table layout for mobile */
@media (max-width: 640px) {
  table,
  thead,
  tbody,
  th,
  td,
  tr {
    display: block;
  }

  thead {
    display: none; /* Hide headers on mobile */
  }

  tr {
    border: 1px solid #ddd;
    margin-bottom: 0.75rem;
    border-radius: 8px;
    overflow: hidden;
  }

  td {
    display: flex;
    justify-content: space-between;
    padding: 0.5rem;
    border-bottom: 1px solid #ddd;
  }

  td:last-child {
    border-bottom: 0;
  }

  td::before {
    content: attr(data-label);
    font-weight: bold;
    color: #374151;
  }
}

</style>
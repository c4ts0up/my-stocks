<template>
  <div class="p-4">
    <h1 class="text-3xl font-bold mb-4 text-green-600 text-center title-font">Stocks List</h1>
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
      <tr v-for="stock in stocks" :key="stock.ticker" class="cursor-pointer hover:bg-gray-100" @click="selectStock(stock.ticker)">
        <td class="p-2 border-b border-gray-300 text-left">{{ stock.ticker }}</td>
        <td class="p-2 border-b border-gray-300 text-left">{{ stock.company_name }}</td>
        <td class="p-2 border-b border-gray-300 text-right">
          <span class="float-left">USD</span>
          <span class="float-right">{{ stock.last_price.toFixed(2) }}</span>
        </td>
        <td class="p-2 border-b border-gray-300 text-center">
          <span :class="getRecommendationClass(stock.recommendation)">{{ stock.recommendation }}</span>
        </td>
      </tr>
      </tbody>
    </table>

    <div v-if="selectedStock" class="fixed inset-0 bg-gray-900 bg-opacity-75 flex items-center justify-center">
      <div class="bg-white p-4 rounded shadow-lg w-96">
        <h2 class="text-xl font-bold mb-2">{{ selectedStock.stock_base.company_name }}</h2>
        <p><strong>Ticker:</strong> {{ selectedStock.stock_base.ticker }}</p>
        <p><strong>Price:</strong> ${{ selectedStock.stock_base.last_price.toFixed(2) }}</p>
        <p><strong>Recommendation:</strong> {{ selectedStock.stock_base.recommendation }}</p>

        <h3 class="mt-4 font-semibold">Ratings:</h3>
        <ul>
          <li v-for="rating in selectedStock.stock_ratings" :key="rating.time" class="text-sm">
            {{ rating.brokerage }}: {{ rating.action }} from {{ rating.rating_from }} to {{ rating.rating_to }} ({{ rating.time }})
          </li>
        </ul>

        <button @click="selectedStock = null" class="mt-4 px-4 py-2 bg-red-500 text-white rounded">Close</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import axios from 'axios';

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

const getRecommendationClass = (recommendation: string) => {
  switch (recommendation) {
    case 'Buy':
      return 'bg-green-500 text-white px-2 py-1 rounded-full';
    case 'Hold':
      return 'bg-gray-400 text-white px-2 py-1 rounded-full';
    case 'Sell':
      return 'bg-red-500 text-white px-2 py-1 rounded-full';
    default:
      return '';
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
  text-transform: uppercase;
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

/* Detail modal styling */
div.fixed {
  padding: 1rem;
}

div.bg-white {
  width: 100%;
  max-width: 400px;
}
</style>
<template>
  <div class="p-4">
    <h1 class="text-2xl font-bold mb-4">Stocks List</h1>
    <div v-if="loading" class="text-center">Loading...</div>
    <div v-else-if="error" class="text-red-500">{{ error }}</div>

    <table v-else class="table-auto w-full border-collapse border border-gray-300">
      <thead class="bg-gray-200">
      <tr>
        <th class="p-2 border border-gray-300">Ticker</th>
        <th class="p-2 border border-gray-300">Company Name</th>
        <th class="p-2 border border-gray-300">Current Price</th>
        <th class="p-2 border border-gray-300">Recommendation</th>
      </tr>
      </thead>
      <tbody>
      <tr v-for="stock in stocks" :key="stock.ticker" class="cursor-pointer hover:bg-gray-100" @click="selectStock(stock.ticker)">
        <td class="p-2 border border-gray-300 text-center">{{ stock.ticker }}</td>
        <td class="p-2 border border-gray-300">{{ stock.company_name }}</td>
        <td class="p-2 border border-gray-300 text-right">${{ stock.last_price.toFixed(2) }}</td>
        <td class="p-2 border border-gray-300 text-center">
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
  font-family: 'Arial', sans-serif;
}
</style>
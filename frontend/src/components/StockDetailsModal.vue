<template>
  <div v-if="stock" class="fixed inset-0 bg-gray-900 bg-opacity-75 flex items-center justify-center">
    <div class="bg-white p-4 rounded shadow-lg w-96">
      <h2 class="text-xl font-bold mb-2">{{ stock.stock_base.company_name }}</h2>
      <p><strong>Ticker:</strong> {{ stock.stock_base.ticker }}</p>
      <p><strong>Price:</strong> ${{ stock.stock_base.last_price.toFixed(2) }}</p>
      <p><strong>Recommendation:</strong> {{ stock.stock_base.recommendation }}</p>

      <h3 class="mt-4 font-semibold">Ratings:</h3>
      <ul>
        <li v-for="rating in stock.stock_ratings" :key="rating.time" class="text-sm">
          {{ rating.brokerage }}: {{ rating.action }} from {{ rating.rating_from }} to {{ rating.rating_to }} ({{ rating.time }})
        </li>
      </ul>

      <button @click="$emit('close')" class="mt-4 px-4 py-2 bg-red-500 text-white rounded">Close</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { defineProps, defineEmits } from 'vue';

const props = defineProps<{ stock: any | null }>();
const emit = defineEmits(['close']);
</script>

<style scoped>
</style>

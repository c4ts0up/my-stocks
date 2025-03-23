<template>
  <div v-if="stock" class="fixed inset-0 bg-gray-300 bg-opacity-75 flex items-center justify-center p-4">
    <div class="bg-white p-6 rounded-lg shadow-xl max-w-[800px] w-full border border-gray-300 text-center">
      <h3 class="text-gray-800">{{ stock.stock_base.ticker }}</h3>
      <h2 class="text-2xl font-semibold mb-3 text-gray-900 title-font">{{ stock.stock_base.company_name }}</h2>
      <h3 class="text-gray-800">USD {{ stock.stock_base.last_price.toFixed(2) }}</h3>


      <span :class="getRecommendationClass(stock.stock_base.recommendation)">
        {{ stock.stock_base.recommendation === "N/A" ? "" : stock.stock_base.recommendation }}
      </span>


      <h3 class="mt-4 font-semibold text-gray-900">Ratings</h3>
      <table>
        <tbody v-for="rating in stock.stock_ratings" :key="rating.time" class="text-sm text-gray-700">
          <tr>
            <td>{{ formatDate(rating.time) }}</td>
            <td>{{ rating.brokerage }}</td>
            <td>{{ rating.rating_from }} → {{ rating.rating_to }}</td>
          </tr>
        </tbody>
      </table>

      <div class="mt-6 flex justify-center">
        <button @click="$emit('close')" class="px-4 py-2 bg-red-500 hover:bg-red-600 text-white font-medium rounded-lg transition">Close</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { defineProps, defineEmits } from 'vue';
import {getRecommendationClass} from "@/utils/recommendationUtils";

const props = defineProps<{ stock: any | null }>();
const emit = defineEmits(['close']);

const formatDate = (rfc3339NanoString: string): string => {
  const date = new Date(rfc3339NanoString);
  return date.toISOString().replace('T', ' ').split('.')[0];
};
</script>

<style scoped>
body {
  font-family: 'Inter', sans-serif;
  background-color: #f9fafb;
}

h2.title-font {
  font-family: 'Poppins', sans-serif;
  letter-spacing: 1px;
  font-size: clamp(1.5rem, 4vw, 2.5rem);
}

h3 {
  font-family: 'Inter', sans-serif;
  letter-spacing: 1px;
  font-size: clamp(1rem, 2vw, 2rem);
}

</style>

// StockRow.vue
<template>
  <tr class="cursor-pointer hover:bg-gray-100" @click="$emit('select', stock.ticker)">
    <td class="p-2 border-b border-gray-300 text-left">{{ stock.ticker }}</td>
    <td class="p-2 border-b border-gray-300 text-left">{{ stock.company_name }}</td>
    <td class="p-2 border-b border-gray-300 text-right">
      <span class="float-left">USD</span>
      <span class="float-right">{{ stock.last_price.toFixed(2) }}</span>
    </td>
    <td class="p-2 border-b border-gray-300 text-center">
      <span :class="getRecommendationClass(stock.recommendation, SIZE.md)">
        {{ stock.recommendation === "N/A" ? "" : stock.recommendation }}
      </span>
    </td>
  </tr>
</template>

<script setup>
import {getRecommendationClass, SIZE} from "@/utils/recommendationUtils.js";

const props = defineProps({
  stock: Object,
});

const emit = defineEmits(["select"]);

</script>
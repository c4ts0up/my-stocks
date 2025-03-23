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
      <span :class="getRecommendationClass(stock.recommendation)">
        {{ stock.recommendation === "N/A" ? "" : stock.recommendation }}
      </span>
    </td>
  </tr>
</template>

<script setup>
const props = defineProps({
  stock: Object,
});

const emit = defineEmits(["select"]);

const getRecommendationClass = (recommendation) => {
  switch (recommendation) {
    case "Buy":
      return "bg-green-500 text-white px-2 py-1 rounded-full";
    case "Hold":
      return "bg-gray-400 text-white px-2 py-1 rounded-full";
    case "Sell":
      return "bg-red-500 text-white px-2 py-1 rounded-full";
    default:
      return "";
  }
};
</script>
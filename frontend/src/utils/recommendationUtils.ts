export const getRecommendationClass = (recommendation: string) => {
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
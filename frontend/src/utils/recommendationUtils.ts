
export enum SIZE {
    sm = "sm",
    md = "md",
    lg = "lg",
    xl = "xl"
}

export const getRecommendationClass = (recommendation: string, size: SIZE = SIZE.xl) => {
    return recommendationTagColor(recommendation) + " rounded-full " + recommendationTagSize(size);
};

const recommendationTagColor = (recommendation: string) => {
    switch (recommendation) {
        case "Buy":
            return "bg-green-500 text-white"
        case "Hold":
            return "bg-gray-400 text-white";
        case "Sell":
            return "bg-red-500 text-white";
        default:
            return "";
    }
}

const recommendationTagSize = (size: SIZE) => {

    let buttonSizes = new Map<SIZE, String>();
    buttonSizes.set(SIZE.sm, "px-3 py-1");
    buttonSizes.set(SIZE.md, "px-5 py-1.5");
    buttonSizes.set(SIZE.lg, "px-5 py-2");
    buttonSizes.set(SIZE.xl, "px-6 py-2");

    return "text-" + size + " " + buttonSizes.get(size);
}
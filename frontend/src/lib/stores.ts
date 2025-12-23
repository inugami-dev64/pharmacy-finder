import { writable } from "svelte/store";
import { PharmacyReview } from "./service/pharmacy-review";

export const reviewData = writable<PharmacyReview[] | undefined>([]);
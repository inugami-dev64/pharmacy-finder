import { writable } from "svelte/store";
import { PharmacyReview } from "./service/pharmacy-review";
import type { PharmacyRating } from "./service/pharmacy-rating";

export const reviewData = writable<PharmacyReview[] | undefined>(undefined);
export const ratingData = writable<PharmacyRating[] | undefined>(undefined);
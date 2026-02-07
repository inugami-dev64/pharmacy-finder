import type { HttpError } from "$lib/http-error";
import { authenticationSession } from "../auth-session";
import { LocalizedBackendError } from "./error";
import { PAGER_LIMIT } from "./pager";
import { PharmacyReview } from "./pharmacy-review";

export enum ReviewModerationResult {
    Approved = "APPROVED",
    PersonalAttack = "PERSONAL_ATTACK",
    Offensive = "OFFENSIVE",
    Other = "OTHER",
    None = "NONE"
}

export class ModeratorPharmacyReview extends PharmacyReview {
    pharmacyId?: number;
    commentReviewResult?: ReviewModerationResult;
    markedForDeletion?: boolean;
    reviewedAt?: number

    /**
     * Queries pharmacy reviews for moderator view.
     * NOTE: Requires authentication
     *
     * @param key optionally specifies the last key value (updatedAt property) of the last page
     * @param uniqueKey optionally specifies the last unique key value (id property) of the last page
     * @param desc specifies whether descending sorting order should be applied (default: false)
     * @param showUnmoderated specifies whether unmoderated reviews should be shown (default: true)
     * @param showModerated specifies whether moderated reviews should be shown (default: false)
     */
    public static async getModeratorPharmacyReviews(
        key?: number,
        uniqueKey?: number,
        desc?: boolean,
        showUnmoderated?: boolean,
        showModerated?: boolean
    ): Promise<ModeratorPharmacyReview[]> {
        const token = authenticationSession.getSessionToken();
        return await fetch(`/api/v1/mod/reviews?l=${PAGER_LIMIT}${key != null ? `&k=${key}` : ""}${uniqueKey != null ? `&uk=${uniqueKey}` : ""}${desc != null ? `&desc=${desc}` : ""}${showModerated != null ? `&moderated=${showModerated}` : ""}${showUnmoderated != null ? `&unmoderated=${showUnmoderated}` : ""}`, {
            method: "GET",
            headers: {
                "Authorization": "Bearer " + token
            }
        })
            .then(async res => {
                if (res.status != 200) {
                    const err: HttpError = await res.json();
                    if (err.msg?.startsWith("Forbidden"))
                        throw new LocalizedBackendError("universalErrors.forbidden");
                    else throw new LocalizedBackendError("universalErrors.internalServerError");
                }

                return await res.json();
            })
    }
}
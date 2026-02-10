import type { HttpError } from "$lib/http-error";
import { authenticationSession } from "../auth-session";
import { LocalizedBackendError } from "./error";

export enum ReviewModerationResult {
    Approved = "APPROVED",
    PersonalAttack = "PERSONAL_ATTACK",
    Offensive = "OFFENSIVE",
    Other = "OTHER",
    None = "NONE"
}

export class Moderation {
    commentId?: number;
    id?: number;
    markedForDeletion?: boolean;
    moderatorComment?: string;
    moderatorId?: string;
    moderatorUsername?: string;
    result?: ReviewModerationResult;
    reviewedAt?: number;

    /**
     * Retrieves all moderations for specified review
     *
     * @param reviewId specifies the review whose moderations to query for
     * @returns a promise to Moderation array
     */
    public static async getModerationsForReview(reviewId: number): Promise<Moderation[]> {
        const token = authenticationSession.getSessionToken();
        return await fetch(`/api/v1/mod/reviews/${reviewId}/moderations`, {
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
                    else if (err.msg?.startsWith("Malformed ID"))
                        throw new LocalizedBackendError("mod.moderations.errors.invalidReviewId");
                    else throw new LocalizedBackendError("universalErrors.internalServerError");
                }

                return await res.json();
            })
    }

    /**
     * Creates a new moderation from given moderation object to specified review
     *
     * @param reviewId specifies the ID of the review to whom to create the review
     */
    public async createModeration(reviewId: number): Promise<void> {
        const token = authenticationSession.getSessionToken();
        const creationData: Moderation = new Moderation;
        creationData.moderatorComment = this.moderatorComment;
        creationData.markedForDeletion = this.markedForDeletion;
        creationData.result = this.result;

        return await fetch(`/api/v1/mod/reviews/${reviewId}/moderations`, {
            method: "POST",
            body: JSON.stringify(creationData),
            headers: {
                "Authorization": "Bearer " + token,
                "Content-Type": "application/json"
            }
        })
            .then(async res => {
                if (res.status != 201) {
                    const err: HttpError = await res.json();
                    if (err.msg?.startsWith("Forbidden"))
                        throw new LocalizedBackendError("universalErrors.forbidden");
                    else if (err.msg?.startsWith("Malformed ID"))
                        throw new LocalizedBackendError("mod.moderations.errors.invalidReviewId");
                    else if (err.msg?.startsWith("Moderation review by this user for comment"))
                        throw new LocalizedBackendError("mod.moderations.errors.moderationExists");
                    else throw new LocalizedBackendError("universalErrors.internalServerError");
                }

                const data: Moderation = await res.json();
                Object.assign(this, data);
            })
    }

    /**
     * Updates the current moderation
     */
    public async updateModeration(): Promise<void> {
        const token = authenticationSession.getSessionToken();
        const modData: Moderation = new Moderation;
        modData.moderatorComment = this.moderatorComment;
        modData.markedForDeletion = this.markedForDeletion;
        modData.result = this.result;

        return await fetch(`/api/v1/mod/reviews/${this.commentId}/moderations/${this.id}`, {
            method: "PATCH",
            body: JSON.stringify(modData),
            headers: {
                "Authorization": "Bearer " + token,
                "Content-Type": "application/json"
            }
        })
            .then(async res => {
                if (res.status != 200) {
                    const err: HttpError = await res.json();
                    if (err.msg?.startsWith("Forbidden") || err.msg?.startsWith("Permission denied"))
                        throw new LocalizedBackendError("universalErrors.forbidden");
                    else if (err.msg?.startsWith("Malformed reviewId"))
                        throw new LocalizedBackendError("mod.moderations.errors.invalidReviewId");
                    else if (err.msg?.startsWith("Malformed modId"))
                        throw new LocalizedBackendError("mod.moderations.errors.invalidModId");
                    else if (err.msg?.startsWith("Not found"))
                        throw new LocalizedBackendError("mod.moderations.errors.moderationNotFound");
                    else throw new LocalizedBackendError("universalErrors.internalServerError");
                }

                const data: Moderation = await res.json();
                Object.assign(this, data);
            })
    }

    /**
     * Deletes the current moderation
     */
    public async deleteModeration(): Promise<void> {
        const token = authenticationSession.getSessionToken();

        return await fetch(`/api/v1/mod/reviews/${this.commentId}/moderations/${this.id}`, {
            method: "DELETE",
            headers: {
                "Authorization": "Bearer " + token
            }
        })
            .then(async res => {
                if (res.status != 200) {
                    const err: HttpError = await res.json();
                    if (err.msg?.startsWith("Forbidden") || err.msg?.startsWith("Permission denied"))
                        throw new LocalizedBackendError("universalErrors.forbidden");
                    else if (err.msg?.startsWith("Malformed reviewId"))
                        throw new LocalizedBackendError("mod.moderations.errors.invalidReviewId");
                    else if (err.msg?.startsWith("Malformed modId"))
                        throw new LocalizedBackendError("mod.moderations.errors.invalidModId");
                    else if (err.msg?.startsWith("Not found"))
                        throw new LocalizedBackendError("mod.moderations.errors.moderationNotFound");
                    else throw new LocalizedBackendError("universalErrors.internalServerError");
                }

                const data: Moderation = await res.json();
                Object.assign(this, data);
            })
    }
}


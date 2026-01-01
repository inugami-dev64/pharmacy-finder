/**
 * Simple data type to represent a geographic point
 */
export class Point {
    lat: number;
    lng: number;

    constructor(lat: number = 0, lng: number = 0) {
        this.lat = lat;
        this.lng = lng;
    }
}
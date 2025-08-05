<script lang="ts">
    import type { PharmacyInfo } from '$lib/service/pharmacy-info';
    import type { LatLngExpression } from 'leaflet';
    import { onMount, onDestroy } from 'svelte';

    export let pharmacies: PharmacyInfo[];
    export let callback: (pharmacy: PharmacyInfo) => Promise<void>;


    let mapElement: HTMLElement;
    let map: any;

    onMount(async () => {
        // Some specific map values
        const MAP_CENTER: LatLngExpression = [58.822, 25.472];
        const MAP_BOUNDS = [
            [57.45, 21.149],
            [60.245, 28.575]
        ];

        const leaflet = await import('leaflet');
        map = leaflet.map(mapElement, { zoomControl: false }).setView(MAP_CENTER, 13);
        map.fitBounds(MAP_BOUNDS)
        leaflet.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
            attribution: '© <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
        }).addTo(map);

        for (let pharmacy of pharmacies) {
            leaflet.marker([<number>pharmacy.latitude, <number>pharmacy.longitude]).
                addTo(map).
                on('click', e => {
                    callback(pharmacy)
                })
        }
    });

    onDestroy(async () => {
        if(map) {
            console.log('Unloading Leaflet map.');
            map.remove();
        }
    });
</script>

<div bind:this={mapElement}></div>

<style>
    div {
        width: 100vw;
        height: 100vh;
    }
</style>


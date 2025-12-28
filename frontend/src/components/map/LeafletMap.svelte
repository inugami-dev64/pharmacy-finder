<script lang="ts">
    import type { PharmacyInfo } from '$lib/service/pharmacy-info';
    import type { Icon, LatLngExpression } from 'leaflet';
    import { onMount, onDestroy } from 'svelte';

    // Pharmacy icons
    import ApothekaMarker from "$lib/assets/markers/apotheka.png";
    import SudameapteekMarker from "$lib/assets/markers/sudameapteek.png";
    import BenuMarker from "$lib/assets/markers/benu.png";
    import EuroapteekMarker from "$lib/assets/markers/euroapteek.png";

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

        const markers: Map<String, Icon> = new Map<String, Icon>([
            ["apotheka", leaflet.icon({iconUrl: ApothekaMarker, iconSize: [32, 32], iconAnchor: [16, 16]})],
            ["südameapteek", leaflet.icon({iconUrl: SudameapteekMarker, iconSize: [32, 32], iconAnchor: [16, 16]})],
            ["benu", leaflet.icon({iconUrl: BenuMarker, iconSize: [32, 32], iconAnchor: [16, 16]})],
            ["euroapteek", leaflet.icon({iconUrl: EuroapteekMarker, iconSize: [32, 32], iconAnchor: [16, 16]})]
        ]);

        map = leaflet.map(mapElement, { zoomControl: false }).setView(MAP_CENTER, 13);
        map.fitBounds(MAP_BOUNDS)
        leaflet.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
            attribution: '© <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
        }).addTo(map);

        for (let pharmacy of pharmacies) {
            let icon = markers.get(pharmacy.chain?.toLowerCase() || "none")
            leaflet.marker(
                    [<number>pharmacy.lat, <number>pharmacy.lng],
                    icon === undefined ? undefined : {icon: icon}).
                addTo(map).
                on('click', e => callback(pharmacy));
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
        width: 100%;
        height: 100%;
    }
</style>


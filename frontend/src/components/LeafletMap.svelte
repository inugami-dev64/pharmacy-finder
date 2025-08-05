<script lang="ts">
    import { onMount, onDestroy } from 'svelte';

    let mapElement: HTMLElement;
    let map: any;

    // Some specific map values
    const MAP_CENTER = [58.822, 25.472]
    const MAP_BOUNDS = [
        [57.45, 21.149],
        [60.245, 28.575]
    ]

    onMount(async () => {
        const leaflet = await import('leaflet');
        map = leaflet.map(mapElement).setView([58.822, 25.472], 13);
        map.fitBounds([
            [57.45, 21.149],
            [60.245, 28.575]
        ])
        leaflet.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
            attribution: '© <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
        }).addTo(map);
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


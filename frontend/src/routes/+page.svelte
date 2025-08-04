<script lang="ts">
    import { _ } from 'svelte-i18n';
</script>

<div class="intro-container">
    <img src="banner.svg" alt="banner">
    <h1>{$_("mainPage.title")}</h1>
    <p>{$_("mainPage.introduction")}</p>
    <h2>{$_("mainPage.gotoMap.pretext")}👉️<a href="/map">{$_("mainPage.gotoMap.title")}</a>👈️</h2>
</div>

<footer>
    <b>Made with love by Karmen aka <a href="https://github.com/inugami-dev64">inugami-dev64</a></b>
</footer>

<style>
    :global(body) {
        margin: 0;
        padding: 0;
        width: 100vw;
        height: 100vh;
    }

    .intro-container {
        display: flex;
        flex-direction: column;
        justify-content: center;
        margin: auto;
        width: 70%;
        height: 100%;
    }

    .intro-container > img {
        margin: 0 auto;
        width: 80%;
        height: auto;
    }

    .intro-container > h1 {
        text-align: center;
    }

    .intro-container > h2 {
        text-align: center;
    }

    footer {
        position: fixed;
        bottom: 0;
        left: 0;
        width: 100%;
        margin-bottom: 0.5em;
        display: block;
        text-align: center;
    }
</style>
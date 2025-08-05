<script lang="ts">
    import { _ } from 'svelte-i18n';
    import TransHeart from '../components/main/TransHeart.svelte'
    import LanguageSwitcher from '../components/common/LanguageSwitcher.svelte';
</script>

<div class="wrapper">
    <LanguageSwitcher/>
    <div class="intro-container">
        <img src="banner.svg" alt="banner" class="banner">
        <h1>{$_("mainPage.title")}</h1>
        <p>{$_("mainPage.introduction")}</p>
        <h2>{$_("mainPage.gotoMap.pretext")}<a href="/map">👉️<span class="rainbow">{$_("mainPage.gotoMap.title")}</span>👈️</a></h2>
    </div>

    <footer>
        <span>Made with <TransHeart width=32 height=32/> by Karmen aka <a href="https://github.com/inugami-dev64">inugami-dev64</a></span>
    </footer>
</div>

<style>
    :global(body) {
        background-color: #fdcece;
    }

    .banner {
        width: 50%;
        max-width: 960px;
        min-width: 300px;
    }

    h2 > a, h2 > a > span {
        font-size: var(--h2-size);
    }

    .wrapper {
        display: flex;
        flex-direction: column;
        min-height: 100vh;
    }

    .rainbow {
        font-weight: 900;
        background: linear-gradient(
            270deg,
            #ff0000,
            #ff9900,
            #ffff00,
            #33cc33,
            #0066ff,
            #6600cc,
            #ff3399,
            #ff0000
        );

        background-size: 200% 200%;
        -webkit-background-clip: text;
        background-clip: text;
        -webkit-text-fill-color: transparent;

        animation: rainbowMove 0.5s linear infinite;
        text-shadow: 0 0 20px rgba(255, 255, 255, 0.1);
        filter: drop-shadow(2px 2px 2px gray)
    }

    @keyframes rainbowMove {
        0% {
            background-position: 0% 50%;
        }
        100% {
            background-position: 100% 50%;
        }
    }

    .intro-container {
        padding-top: 0.5em;
        text-align: center;
        display: flex;
        flex-direction: column;
        justify-content: center;
        margin: auto;
        width: 70%;
        height: 100%;
    }

    @media(max-width: 600px) {
        .intro-container {
            width: 95%;
        }
    }

    .intro-container h2 a {
        white-space: nowrap;
    }

    .intro-container > img {
        margin: 0 auto;
        width: 80%;
        height: auto;
    }

    footer {
        width: 100%;
        margin-bottom: 0.5em;
        display: block;
        text-align: center;
    }
</style>
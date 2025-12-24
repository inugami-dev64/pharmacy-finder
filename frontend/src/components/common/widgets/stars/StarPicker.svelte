<script lang="ts">
    export let name: string;
    export let defaultChecked: number = 5;
</script>

<div class="star-container">
    <input type="radio" name={name} class="star" value="1" required checked={defaultChecked === 1}>
    <input type="radio" name={name} class="star" value="2" required checked={defaultChecked === 2}>
    <input type="radio" name={name} class="star" value="3" required checked={defaultChecked === 3}>
    <input type="radio" name={name} class="star" value="4" required checked={defaultChecked === 4}>
    <input type="radio" name={name} class="star" value="5" required checked={defaultChecked === 5}>
</div>

<style>
    .star-container {
        display: inline-block;
        user-select: none;
        padding: 0;
        margin: 0;
        outline: none;
    }

    .star {
        margin: 0;
        appearance: none;
        width: 24px;
        height: 24px;
        background: url("data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMjQiIGhlaWdodD0iMjQiIHZlcnNpb249IjEuMSIgdmlld0JveD0iMCAwIDguNDY2NjY2NiA4LjQ2NjY2NjYiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+CiAgICA8ZyBmaWxsPSIjZmZiMDAwIiBzdHJva2Utd2lkdGg9Ii4yMzk3NDYiPgogICAgICAgIDxwYXRoIGQ9Im00LjIzMzI3MDMgMC40MjMzMzMwMWMtMC4yNTgwNjc4IDAtMC41MTYwODU1IDAuMTQ5MDUyNi0wLjYxMDE1MyAwLjQ0NjgyNjZsLTAuNTAzMzA2MSAxLjU5MjY3MDZjLTAuMDczNDkxMyAwLjIzMjYyNTctMC4yODk0MzQ3IDAuMzkxMTQ5LTAuNTMzNzY2MyAwLjM5MTE0OWgtMS41MjE2MzQxYy0wLjYxOTcwNDA0IDAtMC44Nzc1MjI1NCAwLjc5MTMzNTgtMC4zNzYzMDgwNCAxLjE1NTIwMDRsMS4zNTUyNzE1IDAuOTgzOTU1M2MwLjIwNDk4MzQgMC4xNDg4MTMyIDAuMjg0MDU3OSAwLjQxNjUyMDQgMC4xOTMwNzM3IDAuNjUyNjk1N2wtMC41ODc2NTg1IDEuNTI1NzYyOGMtMC4yMjg4MjU0IDAuNTk0MDAwOSAwLjQ1NzgzMjEgMS4xMjA2MDI2IDAuOTczMzM5OSAwLjc0NjI3MWwxLjI4MTY5NTktMC45MzEwODUxYzAuMTk2MzE0OS0wLjE0MjU1NTUgMC40NjI1NzU5LTAuMTQyNTU1NSAwLjY1ODg5MDggMGwxLjI4MTY5NTkgMC45MzEwODUxYzAuNTE1NTA2OSAwLjM3NDMzMTYgMS4yMDIxNzM1LTAuMTUyMjcwMSAwLjk3MzMzOS0wLjc0NjI3MWwtMC41ODc2NTg2LTEuNTI1NzYyOGMtMC4wOTEwMTE0LTAuMjM2MTc1My0wLjAxMTQ0NDQtMC41MDM4ODI1IDAuMTkzNTQzNS0wLjY1MjY5NTdsMS4zNTQ4MDE3LTAuOTgzOTU1M2MwLjUwMTIyNDktMC4zNjM4NjQ2IDAuMjQzMzk2Ni0xLjE1NTIwMDQtMC4zNzYzMDcyLTEuMTU1MjAwNGgtMS41MjE2MzRjLTAuMjQ0MzA1MiAwLTAuNDYwMzEyMy0wLjE1ODUyMzMtMC41MzM3NjcxLTAuMzkxMTQ5bC0wLjUwMzMwNjEtMS41OTI2NzA2Yy0wLjA5NDA4MTItMC4yOTc3NzQtMC4zNTIwODUzLTAuNDQ2ODI2Ni0wLjYxMDE1MjktMC40NDY4MjY2eiIvPgogICAgPC9nPgo8L3N2Zz4=");
    }

    .star:hover {
        cursor: pointer;
    }

    .star:checked ~ .star {
        background: url("data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMjQiIGhlaWdodD0iMjQiIHZlcnNpb249IjEuMSIgdmlld0JveD0iMCAwIDguNDY2NjY2NiA4LjQ2NjY2NjYiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+CiAgICA8ZyBmaWxsPSIjY2ZjZmNmIiBzdHJva2Utd2lkdGg9Ii4yMzk3NDYiPgogICAgICAgIDxwYXRoIGQ9Im00LjIzMzI3MDMgMC40MjMzMzMwMWMtMC4yNTgwNjc4IDAtMC41MTYwODU1IDAuMTQ5MDUyNi0wLjYxMDE1MyAwLjQ0NjgyNjZsLTAuNTAzMzA2MSAxLjU5MjY3MDZjLTAuMDczNDkxMyAwLjIzMjYyNTctMC4yODk0MzQ3IDAuMzkxMTQ5LTAuNTMzNzY2MyAwLjM5MTE0OWgtMS41MjE2MzQxYy0wLjYxOTcwNDA0IDAtMC44Nzc1MjI1NCAwLjc5MTMzNTgtMC4zNzYzMDgwNCAxLjE1NTIwMDRsMS4zNTUyNzE1IDAuOTgzOTU1M2MwLjIwNDk4MzQgMC4xNDg4MTMyIDAuMjg0MDU3OSAwLjQxNjUyMDQgMC4xOTMwNzM3IDAuNjUyNjk1N2wtMC41ODc2NTg1IDEuNTI1NzYyOGMtMC4yMjg4MjU0IDAuNTk0MDAwOSAwLjQ1NzgzMjEgMS4xMjA2MDI2IDAuOTczMzM5OSAwLjc0NjI3MWwxLjI4MTY5NTktMC45MzEwODUxYzAuMTk2MzE0OS0wLjE0MjU1NTUgMC40NjI1NzU5LTAuMTQyNTU1NSAwLjY1ODg5MDggMGwxLjI4MTY5NTkgMC45MzEwODUxYzAuNTE1NTA2OSAwLjM3NDMzMTYgMS4yMDIxNzM1LTAuMTUyMjcwMSAwLjk3MzMzOS0wLjc0NjI3MWwtMC41ODc2NTg2LTEuNTI1NzYyOGMtMC4wOTEwMTE0LTAuMjM2MTc1My0wLjAxMTQ0NDQtMC41MDM4ODI1IDAuMTkzNTQzNS0wLjY1MjY5NTdsMS4zNTQ4MDE3LTAuOTgzOTU1M2MwLjUwMTIyNDktMC4zNjM4NjQ2IDAuMjQzMzk2Ni0xLjE1NTIwMDQtMC4zNzYzMDcyLTEuMTU1MjAwNGgtMS41MjE2MzRjLTAuMjQ0MzA1MiAwLTAuNDYwMzEyMy0wLjE1ODUyMzMtMC41MzM3NjcxLTAuMzkxMTQ5bC0wLjUwMzMwNjEtMS41OTI2NzA2Yy0wLjA5NDA4MTItMC4yOTc3NzQtMC4zNTIwODUzLTAuNDQ2ODI2Ni0wLjYxMDE1MjktMC40NDY4MjY2eiIvPgogICAgPC9nPgo8L3N2Zz4=");
    }
</style>
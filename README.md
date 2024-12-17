## To run the backend

1. Ensure you have **make** installed.

2. Navigate to the `backend` folder, make a copy of the `.env.template` file and rename it to `.env`, then add in your own values.

3. Run `make create_local_db`, then run `make migrate_up` to initialize the database (which is MySQL btw).

4. Run `make run`. The server will be available at http://localhost:8080 (assuming you didn't change the `PORT` value in the `.env` file).

5. You can take a look at the `Makefile` for more useful scripts, and run `make list` to list out all the available targets. Hopefully they all work lol.

## To run the frontend

1. Navigate to the `frontend` folder, make a copy of the `.env.template` file and rename it to `.env`, then add in your server's URL.

    - Assuming nothing was changed including the `PORT` value in the backend's `.env` file, it should be http://localhost:8080/api/v1.

2. Run `npm i`.

3. Run `npm run dev`. The client will be available at http://localhost:5173.

## Todo

-   Implement a toast component to display errors to the client
-   Handle a wack error where some posts have the same key after mapping
- I might or might not have forgotten to manually refresh the access token in some operations (or maybe all of them) and am relying solely on the client's silent refresh which is kinda disastrous lmao
-   Write better docs (and the **About** page, or just scrap it entirely lol)

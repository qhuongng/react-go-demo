import useAuthStore from "../stores/auth.store";

export const formatDateString = (dateString: string) => {
    const date = new Date(dateString);

    const year = date.getUTCFullYear();
    const month = String(date.getUTCMonth() + 1).padStart(2, "0"); // months are 0-indexed
    const day = String(date.getUTCDate()).padStart(2, "0");

    const hours = String(date.getUTCHours()).padStart(2, "0");
    const minutes = String(date.getUTCMinutes()).padStart(2, "0");

    const formattedDateTime = `${hours}:${minutes}, ${day}/${month}/${year}`;

    return formattedDateTime;
};

export const checkAuth = async () => {
    const isLoggedIn = localStorage.getItem("isLoggedIn");

    if (isLoggedIn === "true") {
        try {
            const response = await fetch(`${import.meta.env.VITE_API_URL}/auth/refresh`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                credentials: "include",
            });

            if (response.ok) {
                const data = await response.json();
                useAuthStore.setState({ id: data.data.id, accessToken: data.data.accessToken });
            } else {
                const data = await response.json();
                // the refresh token is probably expired
                console.log(data.errors[0].message);
                useAuthStore.setState({ id: 0, accessToken: "" });
                localStorage.setItem("isLoggedIn", "false");
            }
        } catch (error) {
            console.log(error);
        }
    }
};

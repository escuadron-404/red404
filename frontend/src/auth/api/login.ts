export const loginUser = async (email: string, password: string) => {
	const response = await fetch("http://localhost:8080/api/login", {
		method: "POST",
		headers: {
			"Content-Type": "application/json",
		},
		body: JSON.stringify({ email, password }),
	});
	if (!response.ok) {
		// Try to get error message from backend, or fallback
		let message = "Login failed";
		try {
			const errData = await response.json();
			message = errData.message || message;
		} catch {
			/* ignore JSON parse errors */
		}
		throw new Error(message);
	}
	const data = await response.json();
	return data.token;
};

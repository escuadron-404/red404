export const registerUser = async (email: string, password: string) => {
  const response = await fetch("http://localhost:8080/api/register", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ email, password }),
  });
  if (!response.ok) {
    let message = "sign up failed";
    try {
      const errData = await response.json();
      message = errData.message || message;
    } catch {
      /* ignore JSON parse errors */
    }
    throw new Error(message);
  }
  return await response.json();
};

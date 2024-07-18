export const apiUrl = "https://positive-muskrat-emerging.ngrok-free.app";

export default {
  async requestLogin(): Promise<void> {
      const response = await fetch(`${apiUrl}/api/line/request-login`, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
        }
      });

      if (!response.ok) {
        throw new Error('Network response was not ok');
      }

      const jsonData = await response.json();

      window.location.href = jsonData.requestURL;

  },
}
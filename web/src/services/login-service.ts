export const apiUrl = "https://polished-pheasant-explicitly.ngrok-free.app/api";

export default {
    async requestLogin(): Promise<void> {
        window.location.href = `${apiUrl}/line/request-login`;
      },
}
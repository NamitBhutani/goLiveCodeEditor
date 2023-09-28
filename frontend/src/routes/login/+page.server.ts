import type { Actions, PageServerLoad, RequestEvent } from './$types';
export const load: PageServerLoad = async () => { };
export const actions: Actions = {
    login: async ({ request }) => {
        const body = await request.formData();
        const username = body.get('username');
        const password = body.get('password');
        const expirationDate = new Date();
        expirationDate.setHours(expirationDate.getHours() + 1);
        async () => {
            const response = await fetch('/v1/login', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ username, password }),
            },
            )
            if (response.ok) {
                const responseData = await response.json();
                console.log(responseData);
                const expirationDate = new Date();
                expirationDate.setHours(expirationDate.getHours() + 1);
                document.cookie = `token=${JSON.stringify(responseData.accesstoken)}; expires=${expirationDate.toUTCString()}; path=/`;
            } else {
                console.log('Error:', response.statusText);
            }
        }
    }
}
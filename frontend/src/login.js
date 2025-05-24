export function renderLogin(onLogin) {
    document.getElementById('app').innerHTML = `
        <div class="login">
            <h2>Login</h2>
            <input id="username" placeholder="Enter username" autofocus />
            <input id="password" placeholder="Enter passsword" autofocus />
            <button id="loginBtn">Login</button>
        </div>
    `;
    document.getElementById('loginBtn').onclick = () => {
        const username = document.getElementById('username').value.trim();
        const password = document.getElementById('password').value.trim();
        if (username && password) onLogin(username, password);
    };
}
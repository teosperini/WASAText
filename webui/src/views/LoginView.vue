<template>
	<div class="login-container">
		<div class="login-card">
			<h1>Accedi</h1>
			<form @submit.prevent="doLogin">
				<div class="input-group">
					<input
						type="text"
						v-model="username"
						placeholder="Inserisci il tuo username"
						required
					/>
				</div>
				<button type="submit" class="login-button">
					<span>Accedi</span>
				</button>
			</form>
			<p v-if="error" class="error">{{ error }}</p>
		</div>
	</div>
</template>

<script>
import axios from "@/services/axios";

export default {
	name: "LoginView",
	data() {
		return {
			username: "",
			error: null
		};
	},
	methods: {
		async doLogin() {
			try {
				const response = await axios.post("/session", {
					username: this.username
				});

				const data = response.data;

				if (!data.userId || data.userId <= 0) {
					throw new Error("ID utente non valido ricevuto dal server!");
				}

				localStorage.setItem("username", String(this.username));
				localStorage.setItem("token", String(data.userId));
				console.log("Token salvato:", localStorage.getItem("token"));

				this.$router.push("/chat");
			} catch (err) {
				console.error(err);
				const status = err.response?.status;
				if (status === 422) {
					this.error = "Username non valido";
				} else if (status >= 400 && status < 500) {
					this.error = "Richiesta non valida";
				} else {
					this.error = "Errore del server: riprova più tardi";
				}
				alert(this.error)
			}
		}
	}
};
</script>

<style scoped>
.login-container {
	display: flex;
	justify-content: center;
	align-items: center;
	height: 100vh;
	background: linear-gradient(135deg, #667eea, #764ba2);
}

.login-card {
	background: white;
	padding: 30px;
	border-radius: 12px;
	box-shadow: 0 5px 15px rgba(0, 0, 0, 0.2);
	text-align: center;
	max-width: 400px;
	width: 100%;
}

h1 {
	color: #333;
	margin-bottom: 20px;
}

.input-group {
	margin-bottom: 20px;
}

input {
	width: 100%;
	padding: 10px;
	border: 1px solid #ccc;
	border-radius: 6px;
	font-size: 16px;
}

.login-button {
	width: 100%;
	background: #667eea;
	color: white;
	border: none;
	padding: 12px;
	border-radius: 6px;
	font-size: 16px;
	cursor: pointer;
	transition: background 0.3s;
}

.login-button:hover {
	background: #5563c1;
}

.error {
	color: red;
	margin-top: 10px;
}
</style>

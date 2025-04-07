<template>
	<div class="modal-overlay" @click.self="$emit('close')">
		<div class="profile-modal">
			<h2>Il tuo profilo</h2>

			<div class="profile-info">
				<div class="profile-img-wrapper" @click="toggleMenu">
					<img :src="newProfileImageUrl || profileImageUrl" class="profile-img" />
					<div class="image-menu" v-if="showImageMenu">
						<button class="image-menu-view" @click.stop="viewImage">Visualizza immagine</button>
						<button class="image-menu-change" @click.stop="triggerFileInput">Cambia immagine</button>
					</div>
				</div>

				<div class="username-edit">
					<input
						type="text"
						v-model="newUsername"
						placeholder="Modifica nome utente"
					/>
				</div>
			</div>

			<input ref="fileInput" type="file" style="display: none" @change="handleFileUpload" />

			<div class="modal-actions">
				<button class="cancel-btn" @click="$emit('close')">Annulla</button>
				<button class="save-btn" @click="save">Salva</button>
			</div>
		</div>
	</div>
</template>

<script>
export default {
	name: "UserProfileModal",
	props: ["username", "uploadImage", "profileImageUrl"],
	data() {
		return {
			newUsername: this.username,
			newProfileImageUrl: "",
			showImageMenu: false,
		};
	},
	methods: {
		async handleFileUpload(e) {
			const file = e.target.files[0];
			if (!file || !this.uploadImage) return;
			const url = await this.uploadImage(file);
			if (url) this.newProfileImageUrl = url;
			this.showImageMenu = false;
		},
		save() {
			this.$emit("save-profile", {
				newUsername: this.newUsername,
				newProfileImageUrl: this.newProfileImageUrl || null,
			});
			this.$emit("close"); // chiude il modale
		},
		viewImage() {
			window.open(this.newProfileImageUrl || this.profileImageUrl, "_blank");
			this.showImageMenu = false;
		},
		toggleMenu() {
			this.showImageMenu = !this.showImageMenu;
		},
		triggerFileInput() {
			this.$refs.fileInput.click();
			this.showImageMenu = false;
		},
	},
};
</script>


<style scoped>
.modal-overlay {
	position: fixed;
	inset: 0;
	background: rgba(0, 0, 0, 0.4);
	display: flex;
	justify-content: center;
	align-items: center;
	z-index: 200;
	animation: fadeIn 0.25s ease;
}

.profile-modal {
	background: white;
	padding: 30px;
	border-radius: 20px;
	width: 90%;
	max-width: 500px;
	box-shadow: 0 10px 25px rgba(0, 0, 0, 0.15);
	display: flex;
	flex-direction: column;
	gap: 20px;
	animation: slideUp 0.3s ease-out;
}

.profile-info {
	display: flex;
	align-items: center;
	gap: 20px;
}

.profile-img {
	width: 100px;
	height: 100px;
	border-radius: 50%;
	object-fit: cover;
	border: 3px solid #667eea;
	box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
}

.profile-img-wrapper {
	position: relative;
	cursor: pointer;
}

.image-menu {
	position: absolute;
	top: 110%;
	left: 0;
	background: white;
	border: 1px solid #ccc;
	border-radius: 8px;
	box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
	display: flex;
	flex-direction: column;
	z-index: 10;
}

.image-menu button {
	background: none;
	border: none;
	padding: 10px 15px;
	cursor: pointer;
	text-align: left;
	width: 100%;
	transition: background 0.2s;
}

.image-menu button:hover {
	background: #f0f0f0;
	border-radius: 8px;
}

.username-edit input {
	padding: 10px;
	font-size: 1.1rem;
	border-radius: 10px;
	border: 1px solid #ccc;
	flex: 1;
	width: 100%;
}

.upload-section label {
	font-weight: bold;
	margin-bottom: 5px;
}

.upload-section input[type="text"],
.upload-section input[type="file"] {
	padding: 8px;
	border-radius: 8px;
	border: 1px solid #ccc;
	width: 100%;
}

.modal-actions {
	display: flex;
	justify-content: space-between;
	gap: 10px;
}

.save-btn,
.cancel-btn {
	flex: 1;
	padding: 10px 15px;
	border: none;
	border-radius: 10px;
	font-weight: bold;
	color: white;
	cursor: pointer;
}

.save-btn {
	background: linear-gradient(135deg, #667eea, #5563c1);
}

.cancel-btn {
	background: #d9534f;
}

.save-btn:hover {
	background: linear-gradient(135deg, #5563c1, #4452a0);
}

.cancel-btn:hover {
	background: #c9302c;
}

@keyframes fadeIn {
	from {
		background: rgba(0, 0, 0, 0);
	}
	to {
		background: rgba(0, 0, 0, 0.4);
	}
}

@keyframes slideUp {
	from {
		transform: translateY(30px);
		opacity: 0;
	}
	to {
		transform: translateY(0);
		opacity: 1;
	}
}
</style>

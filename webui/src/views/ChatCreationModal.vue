<template>
	<div class="modal-overlay">
		<div class="custom-modal">
			<div class="modal-body">
				<!-- Colonna sinistra -->
				<div class="left-panel">
					<!-- Tipo conversazione -->
					<div class="conversation-type">
						<button
							:class="{ active: chatType === 'private' }"
							@click="chatType = 'private'"
							type="button"
						>
							👤 Privata
						</button>
						<button
							:class="{ active: chatType === 'group' }"
							@click="chatType = 'group'"
							type="button"
						>
							👥 Gruppo
						</button>
					</div>

					<!-- Campi gruppo -->
					<div v-if="chatType === 'group'" class="group-fields">
						<div class="group-row">
							<div class="input-wrapper">
								<label for="groupName">Nome del gruppo:</label>
								<input
									type="text"
									id="groupName"
									v-model="groupName"
									placeholder="Inserisci il nome del gruppo"
								/>
							</div>
							<div class="input-wrapper">
								<label for="chatImageFile">Immagine (opzionale):</label>
								<input type="file" id="chatImageFile" @change="handleFileUpload" />
							</div>

						</div>
					</div>

					<!-- Ricerca utenti -->
					<input
						type="text"
						v-model="searchQuery"
						placeholder="Cerca utente..."
						class="user-search"
					/>
					<p class="select-info">
						Seleziona {{ chatType === 'group' ? 'più utenti' : 'un utente' }}:
					</p>
					<!-- Utenti selezionati -->
					<div v-if="selectedUsersDetails.length" class="selected-users-bar">
						<div
							v-for="user in selectedUsersDetails"
							:key="user.username"
							class="selected-user"
						>
							<img :src="user.profileImageUrl || 'default.png'" class="selected-user-img" />
							<span>{{ user.username }}</span>
							<button class="remove-btn" @click="toggleUser(user.username)">×</button>
						</div>
					</div>
				</div>

				<!-- Colonna destra: lista utenti -->
				<div class="right-panel">
					<ul class="user-list">
						<li
							v-for="user in filteredUsers"
							:key="user.username"
							:class="{ selected: selectedUsers.includes(user.username) }"
							@click="toggleUser(user.username)"
						>
							<img :src="user.profileImageUrl || 'default.png'" class="user-img" />
							<span>{{ user.username }}</span>
						</li>
					</ul>
				</div>
			</div>
			<div class="modal-actions">
				<button class="cancel-btn" @click="$emit('close')">Annulla</button>
				<button class="create-btn" @click="submit">Crea</button>
			</div>
		</div>
	</div>
</template>

<script>
export default {
	name: "ChatCreationModal",
	props: {
		users: Array,
		uploadImage: Function,
	},
	data() {
		return {
			chatType: "private",
			selectedUsers: [],
			groupName: "",
			chatImageUrl: "",
			searchQuery: "",
			imageFile: null,
		};
	},
	watch: {
		chatType() {
			this.selectedUsers = [];
		}
	},
	computed: {
		selectedUsersDetails() {
			return this.users.filter(user => this.selectedUsers.includes(user.username));
		},
		filteredUsers() {
			const currentUsername = localStorage.getItem("username");
			return this.users
				.filter(user => user.username !== currentUsername)
				.filter(user =>
					user.username.toLowerCase().includes(this.searchQuery.toLowerCase())
				);
		}
	},
	methods: {
		async handleFileUpload(event) {
			this.imageFile = event.target.files[0];
		},
		async submit() {
			if (this.selectedUsers.length === 0) {
				alert("Seleziona almeno un utente");
				return;
			}
			if (this.chatType === "group" && !this.groupName.trim()) {
				alert("Inserisci il nome del gruppo");
				return;
			}

			let uploadedUrl = this.chatImageUrl;
			if (this.imageFile && this.uploadImage) {
				uploadedUrl = await this.uploadImage(this.imageFile);
			}

			this.$emit("chat-created", {
				chatType: this.chatType,
				members: this.selectedUsers,
				groupName: this.chatType === "group" ? this.groupName.trim() : null,
				chatImageUrl: uploadedUrl || null
			});
		},
		toggleUser(username) {
			if (this.chatType === "private") {
				if (this.selectedUsers.includes(username)) {
					this.selectedUsers = [];
				} else {
					this.selectedUsers = [username];
				}
			} else {
				if (this.selectedUsers.includes(username)) {
					this.selectedUsers = this.selectedUsers.filter(u => u !== username);
				} else {
					this.selectedUsers.push(username);
				}
			}
		},
	}
};
</script>

<style scoped>
.modal-overlay {
	position: fixed;
	top: 0;
	left: 0;
	width: 100%;
	height: 100%;
	background: rgba(10, 10, 20, 0.8);
	backdrop-filter: blur(6px);
	display: flex;
	justify-content: center;
	align-items: center;
	z-index: 100;
}

.custom-modal {
	background: linear-gradient(135deg, #f7c59f, #a1c4fd);
	color: #1c1c1e;
	padding: 2% 3%;
	border-radius: 20px;
	box-shadow: 0px 0px 15px rgba(0, 0, 0, 0.3);
	width: 85%;
	max-width: 900px;
	overflow-y: auto;
	display: flex;
	flex-direction: column;
	gap: 1.2rem;
	animation: slideUp 0.3s ease-out;
}

.modal-body {
	display: flex;
	gap: 4%;
	width: 100%;
}

.left-panel {
	width: 58%;
	display: flex;
	flex-direction: column;
	gap: 12px;
}

.right-panel {
	width: 38%;
	max-height: 400px;
	overflow-y: auto;
	background: rgba(255, 255, 255, 0.2);
	border-radius: 12px;
	padding: 10px;
}

.conversation-type {
	display: flex;
	gap: 12px;
	justify-content: center;
}

.conversation-type button {
	padding: 10px 20px;
	border: none;
	border-radius: 15px;
	background: #eee;
	color: #333;
	font-weight: bold;
	cursor: pointer;
	transition: all 0.3s;
}

.conversation-type button.active {
	background: #ffd700;
	color: #1c1c1e;
	box-shadow: 0 0 5px rgba(0, 0, 0, 0.2);
}

.group-row {
	display: flex;
	gap: 5%;
	justify-content: space-between;
}

.input-wrapper {
	width: 47.5%;
	display: flex;
	flex-direction: column;
}

.group-fields label {
	font-weight: bold;
	margin-bottom: 4px;
}

.group-fields input {
	padding: 8px;
	border: none;
	border-radius: 10px;
	background: #fdf6e3;
	width: 100%;
	font-size: 1rem;
	box-shadow: inset 1px 1px 3px rgba(0, 0, 0, 0.05);
	color: #333;
}

.user-search {
	padding: 10px;
	border-radius: 10px;
	border: none;
	outline: none;
	width: 100%;
	font-size: 1rem;
	background: #f5f5ff;
	color: #333;
}

.user-list {
	list-style: none;
	padding: 0;
	margin: 0;
}

.user-list li {
	display: flex;
	align-items: center;
	gap: 10px;
	margin-bottom: 8px;
	cursor: pointer;
	padding: 10px;
	border-radius: 12px;
	transition: all 0.3s ease;
	background: rgba(255, 255, 255, 0.4);
	position: relative;
	box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.user-list li:hover {
	transform: scale(1.02);
	background: rgba(255, 255, 255, 0.6);
}

.user-list li.selected {
	background: linear-gradient(135deg, #ffeaa7, #fab1a0);
	border: 2px solid #fdcb6e;
	box-shadow: 0 0 8px rgba(253, 203, 110, 0.7);
	transform: scale(1.03);
	font-weight: bold;
	color: #2d3436;
}


.user-img {
	width: 36px;
	height: 36px;
	border-radius: 50%;
	object-fit: cover;
	border: 2px solid #fff;
}

.selected-users-bar {
	display: flex;
	flex-wrap: wrap;
	max-height: 100px;
	overflow-y: auto;
	gap: 10px;
	padding: 10px;
	background: rgba(255, 255, 255, 0.1);
	border-radius: 10px;
}

.select-info {
	margin-top: 6px;
	font-weight: bold;
}

.selected-user {
	display: flex;
	align-items: center;
	background: rgba(255, 255, 255, 0.3);
	border-radius: 15px;
	padding: 4px 8px;
	gap: 6px;
	white-space: nowrap;
	color: #1c1c1e;
	min-height: 34px;
}

.selected-user-img {
	width: 28px;
	height: 28px;
	border-radius: 50%;
	object-fit: cover;
	border: 1px solid #fff;
}

.remove-btn {
	background: transparent;
	border: none;
	color: #333;
	font-size: 1rem;
	cursor: pointer;
	font-weight: bold;
}
.remove-btn:hover {
	color: #ff5050;
}

.modal-actions {
	display: flex;
	justify-content: space-between;
	gap: 10px;
}

.modal-actions button {
	padding: 10px 20px;
	border: none;
	border-radius: 10px;
	cursor: pointer;
	transition: background 0.2s;
	font-weight: bold;
	font-size: 1rem;
}

.create-btn {
	background: #ffd700;
	color: #333;
}

.create-btn:hover {
	background: #ffcf33;
}

.cancel-btn {
	background: #aaa;
	color: white;
}
.cancel-btn:hover {
	background: #888;
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

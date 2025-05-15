<template>
	<div class="modal-overlay" @click.self="$emit('close')">
		<div class="group-profile-modal">
			<h2>Profilo Gruppo</h2>

			<div class="group-info">
				<div class="group-img-wrapper" @click="toggleMenu">
					<img :src="newImageUrl || mediaUrl" class="group-img" />
					<div class="image-menu" v-if="showImageMenu">
						<button @click.stop="viewImage">Visualizza immagine</button>
						<button @click.stop="triggerFileInput">Cambia immagine</button>
					</div>
				</div>

				<input
					type="text"
					v-model="newName"
					placeholder="Modifica nome gruppo"
					class="group-name-input"
				/>
			</div>

			<input ref="fileInput" type="file" style="display: none" @change="handleFileUpload" />

			<h3>Partecipanti</h3>
			<ul class="members-list">
				<li
					v-for="user in memberDetails"
					:key="user.username"
					@click="$emit('view-user', user)"
				>
					<img :src="user.profileImageUrl || 'default.png'" class="member-img" />
					<span>{{ user.username }}</span>
				</li>
			</ul>

			<div class="group-actions">
				<button class="add-btn" @click="$emit('add-member')">➕ Aggiungi membro</button>
				<button class="leave-btn" @click="$emit('leave-group')">🚪 Esci dal gruppo</button>
			</div>

			<div class="modal-actions">
				<button class="cancel-btn" @click="$emit('close')">Annulla</button>
				<button class="save-btn" @click="save">Salva</button>
			</div>
		</div>
	</div>
</template>

<script>
export default {
	name: "GroupProfileModal",
	props: {
		groupName: String,
		mediaUrl: String,
		members: Array, // array di username
		users: Array,   // array completo di oggetti { username, profileImageUrl }
		uploadImage: Function
	},
	data() {
		return {
			newName: this.groupName,
			newImageUrl: "",
			showImageMenu: false,
		};
	},
	computed: {
		memberDetails() {
			return this.users.filter(u => this.members.includes(u.username));
		}
	},
	mounted() {
		console.log("Media URL ricevuto nel GroupProfileModal:", this.mediaUrl);
		console.log("👥 USERS:", this.users);
		console.log("🧍‍♂️ MEMBERS:", this.members);
	},
	methods: {
		async handleFileUpload(e) {
			const file = e.target.files[0];
			if (!file || !this.uploadImage) return;
			const url = await this.uploadImage(file);
			if (url) this.newImageUrl = url;
			this.showImageMenu = false;
		},
		save() {
			this.$emit("save-group", {
				newGroupName: this.newName,
				newGroupImageUrl: this.newImageUrl || null
			});
			this.$emit("close");
		},
		viewImage() {
			window.open(this.newImageUrl || this.mediaUrl, "_blank");
			this.showImageMenu = false;
		},
		toggleMenu() {
			this.showImageMenu = !this.showImageMenu;
		},
		triggerFileInput() {
			this.$refs.fileInput.click();
			this.showImageMenu = false;
		},
	}
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
}

.group-profile-modal {
	background: white;
	padding: 30px;
	border-radius: 20px;
	width: 90%;
	max-width: 600px;
	box-shadow: 0 10px 25px rgba(0, 0, 0, 0.15);
	display: flex;
	flex-direction: column;
	gap: 20px;
}

.group-info {
	display: flex;
	align-items: center;
	gap: 20px;
}

.group-img {
	width: 100px;
	height: 100px;
	border-radius: 50%;
	object-fit: cover;
	border: 3px solid #667eea;
	box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
}

.group-img-wrapper {
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

.group-name-input {
	flex: 1;
	font-size: 1.2rem;
	padding: 10px;
	border-radius: 10px;
	border: 1px solid #ccc;
	width: 100%;
}

.members-list {
	list-style: none;
	padding: 0;
	margin: 0;
	display: flex;
	flex-direction: column;
	gap: 10px;
	max-height: 200px;
	overflow-y: auto;
}

.members-list li {
	display: flex;
	align-items: center;
	gap: 10px;
	cursor: pointer;
	padding: 10px;
	border-radius: 10px;
	background: #f7f7f7;
	transition: background 0.2s;
}

.members-list li:hover {
	background: #eaeaea;
}

.member-img {
	width: 40px;
	height: 40px;
	border-radius: 50%;
	object-fit: cover;
	border: 2px solid #aaa;
}

.group-actions {
	display: flex;
	justify-content: space-between;
}

.add-btn,
.leave-btn {
	padding: 10px;
	flex: 1;
	border: none;
	border-radius: 10px;
	font-weight: bold;
	cursor: pointer;
	color: white;
}

.add-btn {
	background: #28a745;
	margin-right: 10px;
}
.add-btn:hover {
	background: #218838;
}

.leave-btn {
	background: #d9534f;
}
.leave-btn:hover {
	background: #c9302c;
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
.save-btn:hover {
	background: linear-gradient(135deg, #5563c1, #4452a0);
}

.cancel-btn {
	background: #aaa;
}
.cancel-btn:hover {
	background: #888;
}
</style>

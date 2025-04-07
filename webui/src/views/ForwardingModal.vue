<template>
	<div class="modal-overlay">
		<div class="custom-modal">
			<div class="modal-body">
				<!-- Colonna sinistra -->
				<div class="left-panel">
					<h3>Inoltra messaggio</h3>

					<button class="create-btn rounded-btn" @click="forwardToNewChat">➕ Nuova chat</button>

					<input type="text" v-model="searchQuery" placeholder="Cerca chat..." class="user-search" />

					<p class="select-info">Seleziona una chat:</p>

					<div v-if="selectedChat" class="selected-users-bar">
						<div class="selected-user">
							<img :src="selectedChat.chatImageUrl || 'default.png'" class="selected-user-img" />
							<span>{{ selectedChat.chatName }}</span>
							<button class="remove-btn" @click="selectedChat = null">×</button>
						</div>
					</div>
				</div>

				<!-- Colonna destra -->
				<div class="right-panel">
					<ul class="user-list">
						<li
							v-for="chat in filteredConversations"
							:key="chat.conversationId"
							:class="{ selected: selectedChat?.conversationId === chat.conversationId }"
							@click="selectedChat = chat"
						>
							<img :src="chat.chatImageUrl || 'default.png'" class="user-img" />
							<span>{{ chat.chatName }}</span>
						</li>
					</ul>
				</div>
			</div>

			<div class="modal-actions">
				<button class="cancel-btn" @click="$emit('close')">Annulla</button>
				<button class="create-btn" @click="submitForwarding">Inoltra</button>
			</div>
		</div>
	</div>
</template>

<script>
export default {
	props: {
		conversations: Array,
	},
	data() {
		return {
			selectedChat: null,
			searchQuery: "",
		};
	},
	computed: {
		filteredConversations() {
			return this.conversations.filter(chat =>
				chat.chatName.toLowerCase().includes(this.searchQuery.toLowerCase())
			);
		},
	},
	methods: {
		submitForwarding() {
			if (!this.selectedChat) {
				alert("Seleziona una chat!");
				return;
			}
			this.$emit("forward-to-chat", this.selectedChat.conversationId);
		},
		forwardToNewChat() {
			this.$emit("forward-to-new-chat");
		},
	},
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

.select-info {
	font-weight: bold;
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
	gap: 10px;
	padding: 10px;
	background: rgba(255, 255, 255, 0.1);
	border-radius: 10px;
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

.create-btn.rounded-btn {
	border-radius: 8px;
	padding: 10px 18px;
	font-size: 1rem;
	font-weight: bold;
	box-shadow: 0 2px 6px rgba(0, 0, 0, 0.1);
	transition: all 0.2s ease;
}

.create-btn.rounded-btn:hover {
	transform: scale(1.03);
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

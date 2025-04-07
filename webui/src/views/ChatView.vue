<template>
	<div class="chat-container">
		<!-- Sidebar a sinistra -->
		<div class="sidebar">
			<div class="static-sidebar">
				<header>
					<button class="logout-btn" @click="logout">Logout</button>
					<h2>Chat</h2>
					<button class="profile-btn" @click="showProfileModal = true">
						<img :src="profileImageUrl || 'default.png'" alt="Profilo" class="profile-img-preview" />
						<span>Profilo</span>
					</button>

				</header>
				<div class="search-bar-container">
					<input
						type="text"
						class="search-bar"
						v-model="searchQuery"
						placeholder="Cerca chat..."
					/>
					<button
						class="search-toggle"
						:disabled="!searchQuery"
						@click="searchQuery = ''"
						:class="{ disabled: !searchQuery }"
					>
						❌
					</button>

				</div>

			</div>
			<ul>
				<li v-for="conv in filteredConversations" :key="conv.conversationId" @click="selectConversation(conv)">
					<img :src="getChatImageUrl(conv)" alt="Chat" class="chat-img" />
					<div class="chat-info">
						<h3>{{ conv.chatName }}</h3>
						<p class="last-message" :title="formatLastMessagePreview(conv.lastMessage, conv.chatType)">
							{{ formatLastMessagePreview(conv.lastMessage, conv.chatType) }}</p>
						<p>{{ formatTimestamp(conv.lastMessage.timestamp) }}</p>
						<span v-if="conv.unreadMessages > 0" class="unread-badge">
						  {{ conv.unreadMessages > 99 ? '99+' : conv.unreadMessages }}
						</span>
					</div>
				</li>
			</ul>
			<button class="new-chat-btn" @click="openNewChatModal">➕ Nuova Chat</button>
		</div>

		<!-- Modali -->
		<ChatCreationModal
			v-if="showNewChatModal"
			:users="users"
			@chat-created="handleChatCreated"
			@close="showNewChatModal = false"
			:uploadImage="uploadImage"
		/>

		<UserProfileModal
			v-if="showProfileModal"
			:username="username"
			:profileImageUrl="profileImageUrl"
			@close="showProfileModal = false"
			@save-profile="handleUserSettings"
			:uploadImage="uploadImage"
		/>

		<GroupProfileModal
			v-if="showGroupModal"
			ref="groupProfileModal"
			:groupName="selectedConversation.chatName"
			:mediaUrl="selectedConversation.chatImageUrl"
			:members="selectedConversation.members"
			:uploadImage="uploadImage"
			:users="users"
			@close="showGroupModal = false"
			@save-group="handleGroupSettings"
			@add-member="showAddMemberModal = true; showGroupModal = false"
			@leave-group="handleLeaveGroup"
			@view-user="handleViewUser"
		/>

		<UserProfileReadonlyModal
			v-if="viewingUser"
			:user="viewingUser"
			@close="viewingUser = null"
		/>

		<AddMemberModal
			v-if="showAddMemberModal"
			:users="users"
			:excluded="selectedConversation.members"
			@close="showAddMemberModal = false"
			@members-selected="handleAddMembers"
		/>

		<ConfirmModal
			v-if="showLeaveConfirm"
			title="Sei sicuro di voler uscire dal gruppo?"
			message="Una volta uscito, non riceverai più messaggi da questo gruppo."
			@close="showLeaveConfirm = false"
			@confirm="leaveGroup"
		/>

		<ConfirmModal
			v-if="showDeleteConfirm"
			title="Conferma eliminazione"
			:message="`Sei sicuro di voler eliminare ${selectedForDeletion.length} messaggi${selectedForDeletion.length === 1 ? 'o' : ''}?`"
			@close="showDeleteConfirm = false"
			@confirm="deleteSelectedMessages"
		/>

		<ForwardingModal
			v-if="showForwardingModal"
			:conversations="conversations"
			@close="showForwardingModal = false"
			@forward-to-chat="handleForwardToChat"
			@forward-to-new-chat="handleForwardToNewChat"
		/>

		<EmojiPicker
			:visible="emojiPicker.visible"
			:position="emojiPicker.position"
			@select="handleEmojiSelected"
		/>

		<!-- Chat a destra -->
		<div class="chat-content" v-if="selectedConversation">
			<header class="chat-header">
				<div class="chat-header-left">
					<img :src="getChatImageUrl(selectedConversation)" class="in-chat-img" />
					<h2>{{ selectedConversation.chatName }}</h2>
				</div>
				<div class="chat-header-right">
					<button
						v-if="selectedConversation?.chatType === 'group'"
						@click="openGroupModal"
						class="group-settings-btn"
					>
						⚙️ Impostazioni gruppo
					</button>
					<button @click="toggleDeleteMode" class="delete-mode-btn">
						🗑️
					</button>
				</div>
			</header>
			<div
				v-if="messageMenu.visible"
				class="message-context-menu"
				:style="{ top: `${messageMenu.y}px`, left: `${messageMenu.x}px` }"
			>
				<ul>
					<li @click="startForwarding(messageMenu.message)">🔄 Inoltra</li>
					<li @click="replyToMessage(messageMenu.message)">💬 Rispondi</li>
					<li @click="commentMessage(messageMenu.message)">😀 Reazione</li>
				</ul>
			</div>
			<div class="messages-list">
				<div
					v-for="msg in messages"
					:key="msg.messageId"
					:class="['message', msg.senderUsername === username ? 'sent' : 'received']"
					@contextmenu.prevent="openMessageMenu($event, msg)"
				>

					<template v-if="isDeleteMode && msg.senderUsername === username">
						<input
							type="checkbox"
							class="delete-checkbox"
							:value="msg.messageId"
							v-model="selectedForDeletion"
						/>
					</template>
					<div v-if="msg.isAnswering > 0" class="reply-preview">
						<p class="reply-to">@{{ getRepliedMessageUsername(msg.isAnswering) }}</p>
						<p class="reply-text">{{ getRepliedMessageTextImage(msg.isAnswering) }}</p>
					</div>

					<div class="sender-info" v-if="selectedConversation.chatType === 'group' && msg.senderUsername !== username">
						<img :src="getUserProfileImage(msg.senderUsername)" class="sender-avatar" />
						<p class="sender">{{ msg.senderUsername }}</p>
					</div>

					<p v-if="msg.isForwarded" class="forwarded-label">
						🔁 Messaggio inoltrato
					</p>
					<img
						v-if="msg.messageType === 'image' || msg.messageType === 'text_image'"
						:src="msg.mediaUrl"
						alt="📷 Immagine del messaggio"
						class="message-image"
					/>

					<p v-if="msg.messageType === 'text' || msg.messageType === 'text_image'" class="text">
						{{ msg.text }}
					</p>
					<div class="timestamp-container">
						<!-- Puntini a sinistra (per i messaggi propri) -->
						<button
							v-if="msg.senderUsername === username"
							class="msg-options-btn"
							@click="openMessageMenu($event, msg)"
							title="Opzioni"
						>
							<svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor">
								<circle cx="5" cy="12" r="2" />
								<circle cx="12" cy="12" r="2" />
								<circle cx="19" cy="12" r="2" />
							</svg>
						</button>

						<!-- Emoji per i propri messaggi -->
						<div
							v-if="msg.senderUsername === username && (msg.comments || []).length"
							class="message-reactions"
						>
							<span
								v-for="c in msg.comments || []"
								:key="c.username"
								class="reaction"
								@click="c.username === username && removeEmoji(msg)"
								@mouseenter="hoveredReactions[msg.messageId] = c.username"
								@mouseleave="hoveredReactions[msg.messageId] = null"
							>
								{{ c.emoji }}
								<span
									v-if="hoveredReactions[msg.messageId] === c.username"
									class="reaction-tooltip"
								>
									<img :src="getUserProfileImage(c.username)" class="tooltip-avatar" />
									<span class="tooltip-name">{{ c.username }}</span>
								</span>
							</span>
						</div>

						<!-- Timestamp + spunte (mostrato solo una volta per tutti) -->
						<p class="timestamp">
							{{ formatTimestamp(msg.timestamp) }}
							<span v-if="msg.senderUsername === username" class="tick">
							  <span v-if="msg.isRead" class="double-tick">✔✔</span>
							  <span v-else>✔</span>
							</span>
						</p>

						<!-- Emoji per i messaggi altrui -->
						<div
							v-if="msg.senderUsername !== username && (msg.comments || []).length"
							class="message-reactions"
						>
							<span
								v-for="c in msg.comments || []"
								:key="c.username"
								class="reaction"
								@click="c.username === username && removeEmoji(msg)"
								@mouseenter="hoveredReactions[msg.messageId] = c.username"
								@mouseleave="hoveredReactions[msg.messageId] = null"
							>
								{{ c.emoji }}
								<span
									v-if="hoveredReactions[msg.messageId] === c.username"
									class="reaction-tooltip"
								>
									<img :src="getUserProfileImage(c.username)" class="tooltip-avatar" />
									<span class="tooltip-name">{{ c.username }}</span>
								</span>
							</span>

						</div>

						<!-- Puntini a destra (per messaggi ricevuti) -->
						<button
							v-if="msg.senderUsername !== username"
							class="msg-options-btn"
							@click="openMessageMenu($event, msg)"
							title="Opzioni"
						>
							<svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor">
								<circle cx="5" cy="12" r="2" />
								<circle cx="12" cy="12" r="2" />
								<circle cx="19" cy="12" r="2" />
							</svg>
						</button>
					</div>

				</div>
				<div ref="messagesEnd"></div>
			</div>
			<div v-if="newImageUrl" class="image-preview-container">
				<img :src="newImageUrl" alt="Anteprima immagine" class="image-preview" />
				<button class="remove-image-btn" @click="removeSelectedImage">❌</button>
			</div>
			<button
				v-if="isDeleteMode && selectedForDeletion.length"
				@click="showDeleteConfirm = true"
				class="delete-confirm-btn"
			>
				Elimina {{ selectedForDeletion.length }} messaggi<span v-if="selectedForDeletion.length === 1">o</span>
			</button>
			<footer>
				<div v-if="messageToReply" class="replying-to-box">
					<p><strong>Rispondendo a {{ messageToReply.senderUsername }}:</strong> {{ messageToReply.text }}</p>
					<button @click="messageToReply = null">❌ Annulla</button>
				</div>
				<form @submit.prevent="sendMessage">
					<input type="text" v-model="newMessage" placeholder="Scrivi un messaggio" />
					<label class="upload-btn">
						📷
						<input type="file" ref="fileInput" @change="handleImageUpload" accept="image/*" hidden />
					</label>
					<button type="submit">Invia</button>
				</form>
			</footer>
		</div>
		<div class="messages-placeholder" v-else>
			<p>Seleziona una conversazione per iniziare a chattare</p>
		</div>
	</div>
</template>

<script>
import ChatCreationModal from '@/views/ChatCreationModal.vue';
import UserProfileModal from "@/views/UserProfileModal.vue";
import GroupProfileModal from "@/views/GroupProfileModal.vue";
import UserProfileReadonlyModal from "@/views/UserProfileReadonlyModal.vue";
import AddMemberModal from "@/views/AddMemberModal.vue";
import ConfirmModal from "@/views/ConfirmModal.vue";
import axios from '@/services/axios.js';
import ForwardingModal from "@/views/ForwardingModal.vue";
import EmojiPicker from "@/views/EmojiPicker.vue";


export default {
	name: "ChatView",
	components: {
		ForwardingModal,
		ConfirmModal,
		AddMemberModal,
		ChatCreationModal,
		UserProfileModal,
		GroupProfileModal,
		UserProfileReadonlyModal,
		EmojiPicker,
	},
	data() {
		return {
			conversations: [],
			selectedConversation: null,
			messages: [],
			newMessage: "",
			token: localStorage.getItem("token"),
			username: localStorage.getItem("username"),
			profileImageUrl: localStorage.getItem("profileImageUrl"),
			showNewChatModal: false,
			users: [],
			fakeChatParticipants: [],
			searchQuery: "",
			isPolling: false,
			isCreatingChat: false,
			showProfileModal: false,
			showGroupModal: false,
			viewingUser: null,
			showAddMemberModal: false,
			showLeaveConfirm: false,
			newImageUrl: null,
			isDeleteMode: false,
			selectedForDeletion: [],
			showDeleteConfirm: false,
			messageToForward: null,
			showForwardingModal: false,
			isForwardingToNewChat: false,
			messageMenu: {
				visible: false,
				x: 0,
				y: 0,
				message: null,
			},
			emojiPicker: {
				visible: false,
				position: { x: 0, y: 0 },
				targetMessage: null,
			},
			repliedMessagesMap: {},
			messageToReply: null,
			hoveredReactions: {},
		};
	},
	created() {
		if (!this.token) {
			console.error("Token non valido, reindirizzo alla login.");
			this.$router.push("/login");
		} else {
			this.fetchConversations();
		}
	},
	mounted() {
		if (!this.token) {
			this.$router.push("/login");
		} else {
			this.fetchConversations();
			this.startPolling();
			this.fetchUsers();
			this.fetchProfileImage();
		}
	},
	beforeUnmount() {
		clearInterval(this.pollingInterval);
	},
	computed: {
		filteredConversations() {
			if (!this.searchQuery.trim()) return this.conversations;
			const query = this.searchQuery.trim().toLowerCase();
			return this.conversations.filter(c =>
				c.chatName.toLowerCase().includes(query)
			);
		}
	},
	methods: {
		getUserProfileImage(username) {
			const user = this.users.find(u => u.username === username);
			return user?.profileImageUrl || 'default.png';
		},
		getChatImageUrl(conv) {
			return conv.chatImageUrl && conv.chatImageUrl.trim() !== '' ? conv.chatImageUrl : 'default.png';
		},
		openMessageMenu(e, msg) {
			if (this.messageMenu.visible) {
				this.closeMessageMenu();
				return;
			}

			const rect = e.currentTarget.getBoundingClientRect();
			const menuWidth = 150;
			const menuHeight = 230;
			const windowWidth = window.innerWidth;
			const windowHeight = window.innerHeight;

			let posX, posY;

			// Controllo orizzontale (se va fuori a sinistra → mostra a destra del bottone, altrimenti a sinistra)
			if (rect.left < menuWidth + 10) {
				posX = rect.right; // mostra a destra del pulsante
			} else if (rect.left + menuWidth > windowWidth) {
				posX = windowWidth - menuWidth - 10;
			} else {
				posX = rect.left; // default
			}

			// Controllo verticale (se va fuori in basso → mostra sopra il bottone, altrimenti sotto)
			if (rect.bottom + menuHeight > windowHeight) {
				posY = rect.top - 120;
			} else {
				posY = rect.bottom;
			}

			this.messageMenu = {
				visible: true,
				x: posX,
				y: posY,
				message: msg
			};

			setTimeout(() => {
				document.addEventListener("click", this.closeMessageMenu);
			}, 0);
		},
		closeMessageMenu() {
			this.messageMenu.visible = false;
			document.removeEventListener('click', this.closeMessageMenu);
		},
		startForwarding(msg) {
			this.messageToForward = msg;
			this.showForwardingModal = true;
			this.closeMessageMenu();
		},
		async handleForwardToChat(conversationId) {
			try {
				await axios.post(`/conversations/${conversationId}/messages/${this.messageToForward.messageId}`);
				this.showForwardingModal = false;
				this.messageToForward = null;

				await this.fetchConversations();
				if (this.selectedConversation?.conversationId === conversationId) {
					await this.fetchMessages(conversationId);
					this.scrollToBottomSmooth();
				}
			} catch (err) {
				console.error("Errore inoltro messaggio:", err);
				alert("Errore inoltro messaggio");
			}
		},
		async handleForwardToNewChat() {
			// Attiva il forwarding
			this.isForwardingToNewChat = true;
			this.showForwardingModal = false;

			await this.openNewChatModal();
		},
		replyToMessage(msg) {
			this.messageToReply = msg;
			this.closeMessageMenu();
		},
		commentMessage(msg) {
			this.closeMessageMenu();
			setTimeout(() => {
				const menuWidth = 200;
				const windowWidth = window.innerWidth;
				const margin = 10;

				const isMine = msg.senderUsername === this.username;
				let posX = this.messageMenu.x;

				if (isMine) {
					// Di default lo mettiamo a sinistra del pulsante
					posX = this.messageMenu.x - menuWidth;

					// Se va troppo a sinistra, lo spostiamo un po’ più dentro
					if (posX < margin) {
						posX = margin;
					}
				} else {
					// Di default lo mettiamo a destra
					if (posX + menuWidth > windowWidth - margin) {
						posX = windowWidth - menuWidth - margin;
					}
				}

				this.emojiPicker = {
					visible: true,
					position: {
						x: posX,
						y: this.messageMenu.y,
					},
					targetMessage: msg,
				};
				document.addEventListener("click", this.closeEmojiPicker);
			}, 0);
		},
		async handleEmojiSelected(emoji) {
			const { targetMessage } = this.emojiPicker;
			if (!targetMessage) return;

			try {
				await axios.put(
					`/conversations/${this.selectedConversation.conversationId}/messages/${targetMessage.messageId}/emoji`,
					{ emoji },
					{
						headers: {
							"Content-Type": "application/merge-patch+json"
						}
					}
				);

				// Aggiorna localmente il commento (1 emoji per utente)
				if (!targetMessage.comments) {
					this.$set(targetMessage, 'comments', []);
				}
				const existing = targetMessage.comments.find(c => c.username === this.username);
				if (existing) {
					existing.emoji = emoji;
				} else {
					targetMessage.comments.push({ username: this.username, emoji });
				}

				const index = this.messages.findIndex(m => m.messageId === targetMessage.messageId);
				if (index !== -1) {
					this.messages.splice(index, 1, { ...targetMessage });
				}

			} catch (err) {
				console.error("Errore nel commentare con emoji:", err);
				alert("Errore durante l'invio dell'emoji");
			}
			this.closeEmojiPicker();
		},
		closeEmojiPicker() {
			this.emojiPicker.visible = false;
			this.emojiPicker.targetMessage = null;
			document.removeEventListener("click", this.closeEmojiPicker);
		},
		async removeEmoji(message) {
			try {
				await axios.delete(
					`/conversations/${this.selectedConversation.conversationId}/messages/${message.messageId}/emoji`
				);
				message.comments = message.comments.filter(c => c.username !== this.username);
			} catch (err) {
				console.error("Errore nella rimozione dell'emoji:", err);
				alert("Errore durante la rimozione dell'emoji");
			}
		},
		toggleDeleteMode() {
			this.isDeleteMode = !this.isDeleteMode;
			this.selectedForDeletion = [];
		},
		async deleteSelectedMessages() {
			const convId = this.selectedConversation.conversationId;

			try {
				for (const msgId of this.selectedForDeletion) {
					console.log("Eliminazione messaggio con ID:", msgId);
					await axios.delete(`/conversations/${convId}/messages/${msgId}`);
				}

				this.messages = this.messages.filter(
					msg => !this.selectedForDeletion.includes(msg.messageId)
				);

				if (this.messages.length === 0) {
					this.conversations = this.conversations.filter(c => c.conversationId !== convId);
					this.selectedConversation = null;
				}

				this.isDeleteMode = false;
				this.selectedForDeletion = [];
				this.showDeleteConfirm = false;

			} catch (err) {
				console.error("Errore eliminazione:", err);
				alert("Errore durante l'eliminazione dei messaggi");
			}
		},
		removeSelectedImage() {
			this.newImageUrl = null;
			this.$refs.fileInput.value = ""; // reset anche lato DOM
		},
		async handleImageUpload(event) {
			const file = event.target.files[0];
			if (!file) return;

			const url = await this.uploadImage(file);
			if (url) {
				this.newImageUrl = url;
			}
		},
		formatLastMessagePreview(lastMessage, chatType) {
			if (!lastMessage) return "";

			const { messageType, senderUsername, text } = lastMessage;
			let preview = "";

			if (messageType === "text") preview = text?.trim();
			else if (messageType === "image") preview = "📷 Immagine";
			else if (messageType === "text_image") preview = `📷 ${text?.trim()}`;
			else preview = "Messaggio";

			if (chatType === "group" && senderUsername && senderUsername !== this.username) {
				return `${senderUsername}: ${preview}`;
			}

			return preview;
		},
		handleLeaveGroup() {
			this.showLeaveConfirm = true;
		},
		async leaveGroup() {
			this.showLeaveConfirm = false;
			this.showGroupModal = false;

			try {
				await axios.delete(`/conversations/${this.selectedConversation.conversationId}/members`);
				this.conversations = this.conversations.filter(
					c => c.conversationId !== this.selectedConversation.conversationId
				);
				this.selectedConversation = null;
				this.messages = [];
			} catch (err) {
				console.error("Errore durante l’uscita dal gruppo", err);
				alert("Errore durante l’uscita dal gruppo");
			}
		},
		async fetchProfileImage() {
			try {
				const { data } = await axios.get("/image");
				this.profileImageUrl = data.profileImageUrl;
				localStorage.setItem("profileImageUrl", data.profileImageUrl);
			} catch (err) {
				console.error("Errore nel fetch dell'immagine profilo:", err);
			}
		},
		startPolling() {
			this.pollingInterval = setInterval(async () => {
				if (this.isPolling || this.isCreatingChat) return;

				this.isPolling = true;

				try {
					await this.fetchConversations();
					if (this.selectedConversation && !this.selectedConversation.conversationId.toString().startsWith("temp-")) {
						await this.fetchMessages(this.selectedConversation.conversationId);
					}
				} catch (err) {
					console.error("Errore durante il polling:", err);
				} finally {
					this.isPolling = false;
				}
			}, 5000);
		},
		formatTimestamp(timestamp) {
			const date = new Date(timestamp);
			return date.toLocaleTimeString('it-IT', { hour: '2-digit', minute: '2-digit' });
		},
		async handleAddMembers(usernames) {
			const convId = this.selectedConversation.conversationId;

			for (const username of usernames) {
				try {
					await axios.put(
						`/conversations/${convId}/members`,
						{ username },
						{
							headers: {
								"Content-Type": "application/merge-patch+json"
							}
						}
					);

					// Aggiorna i membri nel GroupProfileModal
					this.$refs.groupProfileModal?.$emit("members-added", [username]);

					// Aggiorna localmente i membri della conversazione selezionata
					this.selectedConversation.members = [...this.selectedConversation.members, username];

				} catch (err) {
					console.error("Errore aggiunta membro:", err);
				}
			}

			this.showAddMemberModal = false;
		},
		async fetchUsersAndUpdateGroup() {
			await this.fetchUsers();
			await this.fetchConversations();
		},
		async handleUserSettings({ newUsername, newProfileImageUrl }) {
			try {
				// Aggiorna username
				if (newUsername && newUsername !== this.username) {
					await axios.put("/username", { username: newUsername });

					const oldUsername = this.username;
					this.username = newUsername;
					localStorage.setItem("username", newUsername);

					// Aggiorna i messaggi visivamente
					this.messages = this.messages.map(msg => {
						if (msg.senderUsername === oldUsername) {
							return { ...msg, senderUsername: newUsername };
						}
						return msg;
					});

					if (
						this.selectedConversation &&
						this.selectedConversation.lastMessage?.senderUsername === oldUsername
					) {
						this.selectedConversation.lastMessage.senderUsername = newUsername;
					}
				}

				// Aggiorna immagine profilo
				if (newProfileImageUrl && newProfileImageUrl !== this.profileImageUrl) {
					await axios.put("/image", { profileImageUrl: newProfileImageUrl });
					this.profileImageUrl = newProfileImageUrl;
					localStorage.setItem("profileImageUrl", newProfileImageUrl);
				}

				this.showProfileModal = false;
			} catch (err) {
				console.error(err);
				const status = err.response?.status;
				if (status === 409) {
					alert("Username già in uso");
				} else if (status === 422) {
					alert("Username non valido");
				} else if (status >= 400 && status < 500) {
					alert("Richiesta non valida");
				} else {
					alert("Errore del server: Riprova più tardi");
				}
			}
		},
		async handleGroupSettings({ newGroupName, newGroupImageUrl }) {
			console.log("Ricevuti:", newGroupName, newGroupImageUrl);
			if (!this.selectedConversation) return;

			try {
				const convId = this.selectedConversation.conversationId;
				let nameChanged = false;
				let imageChanged = false;

				// Aggiorna nome del gruppo
				if (newGroupName && newGroupName !== this.selectedConversation.chatName) {
					await axios.put(`/conversations/${convId}/name`, {
						groupName: newGroupName
					});
					this.selectedConversation.chatName = newGroupName;
					nameChanged = true;
				}

				// Aggiorna immagine del gruppo
				if (newGroupImageUrl && newGroupImageUrl !== this.selectedConversation.chatImageUrl) {
					await axios.put(`/conversations/${convId}/image`, {
						chatImageUrl: newGroupImageUrl
					});
					this.selectedConversation.chatImageUrl = newGroupImageUrl;
					imageChanged = true;
				}

				const idx = this.conversations.findIndex(
					c => c.conversationId === this.selectedConversation.conversationId
				);
				if (idx !== -1) {
					if (nameChanged) this.conversations[idx].chatName = newGroupName;

					if (imageChanged) this.conversations[idx].chatImageUrl = newGroupImageUrl;
				}

				this.showGroupModal = false;
			} catch (err) {
				console.error("Errore nel salvataggio del gruppo:", err);
				alert("Errore durante il salvataggio delle modifiche del gruppo");
			}
		},
		async fetchConversations() {
			try {
				const { data } = await axios.get("/conversations");
				this.conversations = data.conversations;
				this.sortConversationsByLastMessage();
			} catch (err) {
				console.error("Errore nel recupero delle conversazioni:", err);
			}
		},
		async fetchUsers() {
			try {
				const { data } = await axios.get("/users");
				this.users = data.users;
			} catch (err) {
				console.error("Errore nel recupero utenti:", err);
			}
		},
		async openNewChatModal() {
			await this.fetchUsers();
			this.showNewChatModal = true;
		},
		async openGroupModal() {
			await this.fetchUsers();
			this.showGroupModal = true;
		},
		sortConversationsByLastMessage() {
			this.conversations.sort((a, b) => {
				const timeA = new Date(a.lastMessage?.timestamp || 0).getTime();
				const timeB = new Date(b.lastMessage?.timestamp || 0).getTime();
				return timeB - timeA;
			});
		},
		handleChatCreated({ chatType, members, groupName, chatImageUrl }) {
			if (chatType === "private") {
				const existing = this.conversations.find(conv =>
					conv.chatType === "private" &&
					(members.includes(conv.chatName) || conv.chatName === members[0])
				);
				if (existing) {
					this.selectedConversation = existing;
					this.showNewChatModal = false;
					this.scrollToBottom();
					return;
				}
			}
			this.isCreatingChat = true;
			const tempId = "temp-" + Date.now();
			const chatName = chatType === "group" ? groupName : members[0];
			let image = "";

			if (chatType === "group") {
				image = chatImageUrl || "";
			} else {
				const otherUser = this.users.find(u => u.username === members[0]);
				image = otherUser?.profileImageUrl || "";
			}

			const newChat = {
				conversationId: tempId,
				chatType,
				chatName,
				chatImageUrl: image,
				lastMessage: { text: "Nessun messaggio", timestamp: new Date() },
				unreadMessages: 0
			};

			this.conversations.push(newChat);
			this.sortConversationsByLastMessage();
			this.selectedConversation = newChat;
			this.messages = [];
			this.fakeChatParticipants = [...members];
			this.showNewChatModal = false;

			if (this.isForwardingToNewChat && this.messageToForward) {
				this.isForwardingToNewChat = true;
				this.sendMessage();
			}
			this.scrollToBottom();
		},
		handleViewUser(user) {
			this.viewingUser = user;
		},
		async uploadImage(file) {
			if (!file) return null;

			const formData = new FormData();
			formData.append("file", file);

			try {
				const { data } = await axios.post("/upload", formData);
				return data.mediaUrl;
			} catch (err) {
				console.error("Errore durante l'upload dell'immagine:", err);
				alert("Errore durante il caricamento dell'immagine");
				return null;
			}
		},
		async selectConversation(conv) {
			if (
				this.selectedConversation &&
				this.selectedConversation.conversationId.toString().startsWith("temp-") &&
				this.messages.length === 0
			) {
				// Rimuove la chat temporanea dalla lista
				this.conversations = this.conversations.filter(
					c => c.conversationId !== this.selectedConversation.conversationId
				);
			}

			this.selectedConversation = conv;
			this.messages = [];
			this.selectedConversation.unreadMessages = 0;
			await this.fetchMessages(conv.conversationId);
			this.scrollToBottom()
		},
		getRepliedMessageTextImage(id) {
			const msg = this.repliedMessagesMap[id];
			if (!msg) return "[messaggio eliminato]";
			if (msg.messageType === "image") return "📷 Immagine";
			if (msg.messageType === "text_image") return `📷 ${msg.text}`;
			return msg.text;
		},
		getRepliedMessageUsername(id) {
			return this.repliedMessagesMap[id]?.senderUsername || "Utente sconosciuto";
		},
		async fetchMessages(conversationId) {
			try {
				let oldMessage = this.messages
				const response = await axios.get(`/conversations/${conversationId}`);
				this.messages = response.data.messages
				this.messages.sort((a, b) => new Date(a.timestamp) - new Date(b.timestamp));
				console.log(response.data.messages)
				const oldLast = oldMessage?.[oldMessage.length - 1]?.message_id;
				const newLast = this.messages?.[this.messages.length - 1]?.message_id;

				this.repliedMessagesMap = {};
				this.messages.forEach(m => {
					this.repliedMessagesMap[m.messageId] = m;
				});

				if (oldLast !== newLast) {
					this.scrollToBottom();
				}
			} catch (err) {
				console.error("Errore nel caricamento dei messaggi:", err);
			}
		},
		async sendMessage() {
			let newMsg;

			if (this.isForwardingToNewChat && this.messageToForward) {
				newMsg = {
					forwardFromMessageId: this.messageToForward.messageId,
					senderUsername: this.username,
					timestamp: new Date(),
					messageType: this.messageToForward.messageType,
					text: this.messageToForward.text,
					mediaUrl: this.messageToForward.mediaUrl,
				};
			} else {
				if (!this.newMessage.trim() && !this.newImageUrl) return;

				newMsg = {
					senderUsername: this.username,
					timestamp: new Date(),
					messageType: (this.newMessage && this.newImageUrl) ? "text_image" : this.newImageUrl ? "image" : "text",
					text: this.newMessage,
					mediaUrl: this.newImageUrl,
					replyToMessageId: this.messageToReply ? this.messageToReply.messageId : null
				};
			}

			try {
				if (!Array.isArray(this.messages)) this.messages = [];

				if (this.selectedConversation.conversationId.toString().startsWith("temp-")) {
					const response = await axios.post("/conversations", {
						members: this.fakeChatParticipants,
						chatType: this.selectedConversation.chatType,
						groupName: this.selectedConversation.chatName,
						groupImageUrl: this.selectedConversation.chatImageUrl,
						initialMessage: newMsg
					});

					this.selectedConversation.conversationId = response.data.conversationId;
					if (response.data.chatName) {
						this.selectedConversation.chatName = response.data.chatName;
					}
					if (response.data.chatImageUrl) {
						this.selectedConversation.chatImageUrl = response.data.chatImageUrl;
					}
					this.messages.push({ ...newMsg, messageId: response.data.messages[0].messageId });
				} else {
					const res = await axios.post(`/conversations/${this.selectedConversation.conversationId}`, newMsg);
					if (newMsg.replyToMessageId) {
						const replied = this.messages.find(
							m => m.messageId === newMsg.replyToMessageId
						);
						if (replied) {
							this.repliedMessagesMap[newMsg.replyToMessageId] = replied;
						}
					}
					this.messages.push({
						...newMsg,
						messageId: res.data.messageId,
						isAnswering: newMsg.replyToMessageId,
						comments: [],
					});
				}

				this.scrollToBottomSmooth();
				this.newMessage = "";
				this.newImageUrl = null;
				this.isForwardingToNewChat = false;
				this.messageToForward = null;
				this.$refs.fileInput.value = ""; // resetta input file

				const newLastMessage = {
					text: newMsg.text || "📷 Immagine",
					timestamp: new Date().toISOString(),
					senderUsername: this.username,
					messageType: newMsg.messageType
				};

				const idx = this.conversations.findIndex(
					c => c.conversationId === this.selectedConversation.conversationId
				);
				if (idx !== -1) {
					this.conversations[idx].lastMessage = newLastMessage;
					this.sortConversationsByLastMessage();
				}
				this.messageToReply = null;

				this.isCreatingChat = false;
				await this.fetchConversations();

				const updated = this.conversations.find(
					c => c.conversationId === this.selectedConversation.conversationId
				);
				if (updated) {
					this.selectedConversation = updated;
				}
			} catch (err) {
				console.error("Errore durante l'invio del messaggio:", err);
				alert("Errore durante l'invio del messaggio");
			}
		},
		scrollToBottom() {
			this.$nextTick(() => {
				const el = this.$refs.messagesEnd;
				if (el && el.scrollIntoView) {
					el.scrollIntoView({ behavior: "auto" });
				}
			});
		},
		scrollToBottomSmooth() {
			this.$nextTick(() => {
				const el = this.$refs.messagesEnd;
				if (el && el.scrollIntoView) {
					el.scrollIntoView({ behavior: "smooth" });
				}
			});
		},
		logout() {
			clearInterval(this.pollingInterval);
			this.conversations = [];
			this.selectedConversation = null;
			this.messages = [];
			this.newMessage = "";
			localStorage.removeItem("token");
			localStorage.removeItem("username");
			this.$router.push("/login");
			this.showNewChatModal = false;
			this.users = [];
			this.fakeChatParticipants = [];
			this.searchQuery = "";
		}
	}
};
</script>

<style scoped>
:root {
	font-size: clamp(0.9rem, 1vw + 0.5rem, 1.1rem);
	font-family: "Segoe UI", Roboto, sans-serif;
}

.chat-container {
	font-size: 1rem;
	width: 100%;
	display: flex;
	height: 100vh;
	flex-direction: row;
	background: linear-gradient(135deg, #667eea, #764ba2);
}

.profile-btn {
	display: flex;
	align-items: center;
	justify-content: center;
	gap: 10px;
	padding: 10px 15px;
	background: linear-gradient(135deg, #FFD966, #F4C542);
	color: #4A2F0B;
	border-radius: 10px;
	border: none;
	font-weight: bold;
	cursor: pointer;
	box-shadow: 2px 2px 5px rgba(0, 0, 0, 0.2);
	transition: transform 0.1s ease-in-out, box-shadow 0.2s, background 0.3s;
	font-size: 1rem;
}

.profile-btn:hover {
	transform: scale(1.05);
	box-shadow: 3px 3px 10px rgba(0, 0, 0, 0.3);
	background: linear-gradient(135deg, #FFE58A, #F1B93E);
}

.profile-img-preview {
	width: 32px;
	height: 32px;
	border-radius: 50%;
	object-fit: cover;
	border: 1px solid #aaa;
}

.sidebar {
	width: 25%;
	flex: 0 0 25%;
	background: white;
	padding: 0.5% 0.5% 0.1%;
	box-shadow: 2px 0px 10px rgba(0, 0, 0, 0.1);
	border-right: 1px solid #ddd;
	display: flex;
	flex-direction: column;
	height: 100vh;
}

.sidebar header {
	display: flex;
	justify-content: space-between;
	align-items: center;
	margin-bottom: 20px;
}

.sidebar header h2 {
	font-size: 2rem;
	font-weight: bold;
	margin: 0;
}

/* Stile base di tutti i bottoni */
button {
	padding: 10px 15px;
	background: linear-gradient(135deg, #667eea, #5563c1);
	color: white;
	border-radius: 10px;
	border: none;
	font-weight: bold;
	cursor: pointer;
	box-shadow: 2px 2px 5px rgba(0, 0, 0, 0.2);
	transition: transform 0.1s ease-in-out, box-shadow 0.2s, background 0.3s;
}

button:hover {
	transform: scale(1.05);
	box-shadow: 3px 3px 10px rgba(0, 0, 0, 0.3);
	background: linear-gradient(1deg, #667eea, #5563c1);
}

/* Solo per il bottone logout */
.logout-btn {
	background: linear-gradient(135deg, #ff5c5c, #e60000);
}

.logout-btn:hover {
	background: linear-gradient(135deg, #e60000 , #ff5c5c);
}

.new-chat-btn {
	position: fixed;
	bottom: 2%;
	align-items: center;
}

.sidebar ul {
	flex: 1;
	overflow-y: auto;
	overflow-x: hidden;
	list-style: none;
	padding: 0 0 80px;
}

.sidebar ul {
	scrollbar-width: auto;
	scrollbar-color: rgba(102, 126, 234, 0.4) transparent;
}

.sidebar li {
	display: flex;
	align-items: center;
	/* padding: 10px; */
	padding-left: 2%;
	padding-right: 2%;
	cursor: pointer;
	border-bottom: 1px solid #eee;
	transition: background 0.2s, transform 0.1s;
	border-radius: 5px;
	overflow: hidden;
}

.sidebar li:hover {
	background: rgba(102, 126, 234, 0.2);
	transform: scale(1.02);
}

.sidebar li.active {
	background: rgba(102, 126, 234, 0.4);
	font-weight: bold;
}

.chat-img {
	width: 20%;
	height: 20%;
	object-fit: cover;
	aspect-ratio: 1 / 1;
	border-radius: 50%;
	margin-right: 10px;
	flex-shrink: 0;
}

.chat-content {
	flex: 1;
	margin-left: 25%;
	flex-direction: column;
	background: #f7f7f7;
	display: flex;
	box-shadow: 0px 0px 10px rgba(0, 0, 0, 0.2);
	border-radius: 10px;
}

.chat-header {
	display: flex;
	justify-content: space-between;
	align-items: center;
	padding: 1.2% 2%;
	background: white;
	box-shadow: 0px 2px 5px rgba(0, 0, 0, 0.1);
	border: 1px solid #ccc;
}

.chat-header-left {
	display: flex;
	align-items: center;
	gap: 16px;
}

.group-settings-btn {
	background: linear-gradient(135deg, #FFD966, #F4C542);
	color: #4A2F0B;
	border: none;
	border-radius: 10px;
	padding: 8px 14px;
	font-weight: bold;
	font-size: 1rem;
	cursor: pointer;
	box-shadow: 2px 2px 5px rgba(0, 0, 0, 0.2);
	transition: transform 0.1s ease-in-out, box-shadow 0.2s, background 0.3s;
}

.group-settings-btn:hover {
	transform: scale(1.05);
	box-shadow: 3px 3px 10px rgba(0, 0, 0, 0.3);
	background: linear-gradient(135deg, #FFE58A, #F1B93E);
}


.chat-header h2 {
	font-size: 1.5rem;
	margin: 0;
}

.in-chat-img {
	width: 48px;
	height: 48px;
	border-radius: 50%;
	object-fit: cover;
	border: 2px solid #ccc;
}


.messages-placeholder {
	margin-left: 25%;
	display: flex;
	align-items: center;
	justify-content: center;
	height: 100vh;
	font-size: 1.2rem;
	text-align: center;
	color: white;
	flex: 1
}

.messages-list {
	display: flex;
	flex: 1;
	padding: 15px;
	overflow-y: auto;
	flex-direction: column;
	background: linear-gradient(135deg, #667eea, #764ba2);
}

.message {
	display: flex;
	flex-direction: column;
	padding: 0.5% 0.8%;
	border-radius: 10px;
	margin-bottom: 10px;
	max-width: 50%;
	overflow: revert;
	word-break: break-word;
}

.message .text {
	font-size: 1.4rem;
	line-height: 1.4;
	margin: 0;
	white-space: normal;
	word-wrap: break-word;
	overflow-wrap: break-word;
}

.message.sent {
	align-self: flex-end;
	border: 1px solid #E1A64E;
	border-top-right-radius: 0;
	background: linear-gradient(135deg, #FFDE7D, #E1A64E);
	color: #4A2F0B;
	text-align: left;
	align-items: flex-end;
}

.message.received {
	align-self: flex-start;
	border-top-left-radius: 0;
	background: linear-gradient(135deg, #A4E4C6, #68A8AD);
	border: 1px solid #68A8AD;
	color: #1D3B40;
	text-align: left;
	align-items: flex-start;
}

.sender{
	font-weight: bold;
	font-size: 0.9rem;
	margin: 0;
}

.message .message-image {
	width: 98%;
	margin: 1%;
	height: auto;
	max-width: 25vw;
	min-width: 80px;
	border-radius: 8px;
}

.timestamp {
	font-size: 0.8rem;
	color: #51585e;
	margin: 0;
}

footer.text{
	flex: 1;
}

footer {
	background: rgba(255, 255, 255, 0.9);
	padding: 12px 15px;
	border-top: 2px solid #ccc;
	display: flex;
	box-shadow: 0px -2px 5px rgba(0, 0, 0, 0.1);
	backdrop-filter: blur(8px);
}

form {
	display: flex;
	flex: 1;
	align-items: center;
	gap: 10px;
}

input[type="text"] {
	flex: 1;
	padding: 10px;
	border-radius: 5px;
	border: 1px solid #ccc;
	box-sizing: border-box;
}

.static-sidebar {
	position: sticky;
	top: 0;
	background: white;
	z-index: 5;
	width: 100%;
	box-sizing: border-box;
	padding-bottom: 3%;
}

.search-bar-container {
	display: flex;
	align-items: center;
	width: 100%;
	margin-top: 10px;
	margin-bottom: 10px;
	border-radius: 10px;
	box-sizing: border-box;
	background: rgba(255, 255, 255, 0.1);
	backdrop-filter: blur(6px);
	padding: 5px;
	position: relative;
	z-index: 1;
	box-shadow: inset 0 0 5px rgba(0, 0, 0, 0.1);
}

.search-bar {
	flex: 1;
	margin-right: 1%;
	padding: 10px;
	border: none;
	border-radius: 8px;
	outline: none;
	font-size: 1rem;
	background: rgba(255, 255, 255, 0.3);
	color: #333;
}

.search-toggle {
	background: transparent;
	border: none;
	font-size: 1.3rem;
	cursor: pointer;
	color: #555;
	transition: transform 0.2s;
	padding: 5px 10px;
}

.search-toggle:hover:enabled {
	transform: scale(1.1);
	color: #222;
}

.search-toggle:disabled {
	cursor: not-allowed;
	transform: none;
	box-shadow: none;
	background: linear-gradient(135deg, #cccccc, #aaaaaa);
	opacity: 0.3;
	color: #aaa;
}


.upload-btn {
	display: inline-flex;
	align-items: center;
	justify-content: center;
	padding: 10px 15px;
	background: linear-gradient(135deg, #FFD966, #F4C542);
	color: #4A2F0B;
	border-radius: 10px;
	cursor: pointer;
	font-weight: bold;
	font-size: 1.2rem;
	box-shadow: 2px 2px 5px rgba(0, 0, 0, 0.2);
	transition: transform 0.2s, box-shadow 0.2s;
}

.upload-btn:hover {
	transform: scale(1.05);
	box-shadow: 3px 3px 10px rgba(0, 0, 0, 0.3);
	background: linear-gradient(135deg, #FFE58A, #F1B93E);
}

.image-preview-container {
	position: relative;
	display: inline-block;
	margin: 10px 20px;
	max-width: 200px;
}

.image-preview {
	width: 100%;
	max-width: 200px;
	border-radius: 10px;
	box-shadow: 0 0 8px rgba(0, 0, 0, 0.2);
}

.remove-image-btn {
	position: absolute;
	top: -8px;
	right: -8px;
	background: #ff5c5c;
	color: white;
	border: none;
	border-radius: 50%;
	width: 24px;
	height: 24px;
	cursor: pointer;
	font-weight: bold;
	display: flex;
	align-items: center;
	justify-content: center;
	box-shadow: 0 0 5px rgba(0, 0, 0, 0.3);
	transition: background 0.2s;
}

.remove-image-btn:hover {
	background: #e60000;
}

.delete-checkbox {
	transform: scale(1.5);
	margin-bottom: 5px;
}

.chat-header-right {
	display: flex;
	align-items: center;
	gap: 10px;
}

.delete-confirm-btn {
	position: absolute;
	margin: 10px auto;
	align-self: center;
	padding: 12px 20px;
	background: linear-gradient(135deg, #ff5c5c, #e60000);
	color: white;
	border-radius: 10px;
	font-weight: bold;
	font-size: 1.1rem;
	border: none;
	box-shadow: 2px 2px 8px rgba(0, 0, 0, 0.3);
	cursor: pointer;
	transition: transform 0.2s, box-shadow 0.2s;
	width: fit-content;
	max-width: 90%;
}

.delete-confirm-btn:hover {
	transform: scale(1.05);
	box-shadow: 4px 4px 12px rgba(0, 0, 0, 0.4);
	background: linear-gradient(135deg, #e60000 , #ff5c5c);
}

.tick {
	font-size: 0.8rem;
	margin-left: 5px;
	display: inline-block;
}

.tick span {
	font-weight: bold;
}

.double-tick {
	color: grey; /* due spunte grigie */
}

.double-tick.read {
	color: #4fc3f7; /* due spunte blu */
}

.message-context-menu {
	min-width: 140px;
	white-space: nowrap;
	overflow: hidden;
	position: fixed;
	background: rgba(255,255,255,0.95);
	border-radius: 8px;
	box-shadow: 0 4px 12px rgba(0,0,0,0.2);
	z-index: 2000;
}

.message-context-menu ul {
	margin: 0;
	padding: 6px 0;
	list-style: none;
}

.message-context-menu li {
	padding: 6px 12px;
	cursor: pointer;
	transition: background 0.2s;
	border-radius: 4px;
}

.message-context-menu li:hover {
	background: rgba(100, 100, 255, 0.2);
}

.forwarded-label {
	font-size: 0.8rem;
	font-style: italic;
	color: #666;
	margin-bottom: 5px;
}

.msg-options-btn {
	background: none;
	border: none;
	padding: 0 4px;
	margin-left: 6px;
	cursor: pointer;
	color: #666;
	opacity: 0;
	transition: opacity 0.2s ease;
	vertical-align: middle;
}

.message:hover .msg-options-btn {
	opacity: 1;
}

.msg-options-btn:hover{
	background: rgba(0, 0, 0, 0.2);
}

.msg-options-btn svg {
	display: inline-block;
	vertical-align: middle;
}

.timestamp-container {
	display: flex;
	align-items: center;
	gap: 6px;
	font-size: 0.8rem;
	color: #51585e;
	margin-top: 5px;
}

.msg-options-btn {
	background: none;
	border: none;
	padding: 0;
	margin: 0;
	cursor: pointer;
	color: #777;
	opacity: 0;
	transition: opacity 0.2s ease;
}

.message:hover .msg-options-btn {
	opacity: 1;
}

.msg-options-btn svg {
	display: block;
}

.reply-preview {
	background: rgba(255, 255, 255, 0.7);
	width: 100%;
	box-sizing: border-box;
	padding: 8px;
	border-left: 4px solid #888;
	margin-bottom: 5px;
	border-radius: 4px;
}
.reply-to {
	font-weight: bold;
	margin: 0;
}
.reply-text {
	margin: 0;
	font-style: italic;
}

.replying-to-box {
	position: relative;
	margin: 0 10px 0 10px;
	padding: 5px 6px;
	background: #fff4cc;
	border-left: 4px solid #f0b400;
	border-radius: 6px;
	font-size: 0.9rem;
	color: #333;
	display: flex;
	align-items: center;
	gap: 10px;
	max-width: 70%;
	box-shadow: 0 2px 6px rgba(0, 0, 0, 0.1);
}
.replying-to-box p {
	margin: 0;
	flex: 1;
	white-space: nowrap;
	overflow: hidden;
	text-overflow: ellipsis;
}

.replying-to-box button {
	background: none;
	border: none;
	color: #d32f2f;
	font-size: 1.1rem;
	cursor: pointer;
}

.message-reactions {
	display: flex;
	gap: 6px;
	margin-top: 4px;
	font-size: 1.2rem;
}

.reaction {
	background: white;
	padding: 2px 6px;
	border-radius: 20px;
	box-shadow: 0 1px 3px rgba(0,0,0,0.2);
	cursor: default;
}

.reaction {
	position: relative;
	background: white;
	padding: 2px 6px;
	border-radius: 20px;
	box-shadow: 0 1px 3px rgba(0,0,0,0.2);
	cursor: default;
}

.reaction-tooltip {
	position: absolute;
	bottom: 120%;
	left: 50%;
	transform: translateX(-50%);
	background: rgba(0, 0, 0, 0.85);
	color: white;
	padding: 6px 10px;
	border-radius: 10px;
	display: inline-flex;
	align-items: center;
	gap: 8px;
	font-size: 0.85rem;
	white-space: nowrap;
	box-shadow: 0 2px 6px rgba(0, 0, 0, 0.3);
	z-index: 10;
}

.tooltip-avatar {
	width: 24px;
	height: 24px;
	border-radius: 50%;
	object-fit: cover;
	border: 1px solid #ccc;
}

.tooltip-name {
	font-weight: bold;
}

.sender-info {
	display: flex;
	align-items: center;
	gap: 6px;
	margin-bottom: 2px;
	font-size: 0.9rem;
}

.sender-avatar {
	width: 28px;
	height: 28px;
	border-radius: 50%;
	object-fit: cover;
	border: 1px solid #ccc;
	box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.last-message {
	white-space: nowrap;
	overflow: hidden;
	text-overflow: ellipsis;
	max-width: 90%;
}

.chat-info {
	flex: 1;
	overflow: hidden;
	display: flex;
	flex-direction: column;
}

.unread-badge {
	position: absolute;
	bottom: 15%;
	right: 2%;
	background: #ff3b30;
	color: white;
	font-size: 0.75rem;
	font-weight: bold;
	padding: 4px 6px;
	border-radius: 50%;
	min-width: 22px;
	height: 22px;
	display: flex;
	align-items: center;
	justify-content: center;
	box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
	z-index: 2;
	transition: transform 0.2s;
}

.sidebar li {
	position: relative;
}

</style>

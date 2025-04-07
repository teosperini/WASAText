<template>
	<div class="modal-overlay" @click.self="$emit('close')">
		<div class="confirm-modal">
			<h2>{{ title }}</h2>
			<p>{{ message }}</p>
			<div class="modal-actions">
				<button class="cancel-btn" @click="$emit('close')">Annulla</button>
				<button class="confirm-btn" @click="$emit('confirm')">Conferma</button>
			</div>
		</div>
	</div>
</template>

<script>
export default {
	props: {
		title: String,
		message: String
	}
};
</script>

<style scoped>
.modal-overlay {
	position: fixed;
	inset: 0;
	background: rgba(0, 0, 0, 0.5);
	display: flex;
	justify-content: center;
	align-items: center;
	z-index: 200;
}

.confirm-modal {
	background: white;
	border-radius: 15px;
	padding: 25px;
	max-width: 400px;
	width: 90%;
	box-shadow: 0 4px 15px rgba(0, 0, 0, 0.2);
	text-align: center;
}

.modal-actions {
	margin-top: 20px;
	display: flex;
	justify-content: space-between;
	gap: 15px;
}

.cancel-btn,
.confirm-btn {
	flex: 1;
	padding: 10px;
	border: none;
	border-radius: 10px;
	font-weight: bold;
	color: white;
	cursor: pointer;
}

.cancel-btn {
	background: #aaa;
}
.cancel-btn:hover {
	background: #888;
}

.confirm-btn {
	background: #d9534f;
}
.confirm-btn:hover {
	background: #c9302c;
}
</style>

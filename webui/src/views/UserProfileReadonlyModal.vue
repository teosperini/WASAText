<template>
	<div class="modal-overlay" @click.self="$emit('close')">
		<div class="readonly-user-modal">
			<img :src="user.profileImageUrl || 'default.png'" class="readonly-user-img" />
			<h2>{{ user.username }}</h2>
			<button class="close-btn" @click="$emit('close')">❌ Chiudi</button>
		</div>
	</div>
</template>

<script>
export default {
	name: "UserProfileReadonlyModal",
	props: {
		user: Object
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
	z-index: 999;
}

.readonly-user-modal {
	background: white;
	padding: 30px;
	border-radius: 20px;
	width: 300px;
	display: flex;
	flex-direction: column;
	align-items: center;
	gap: 20px;
	box-shadow: 0 10px 25px rgba(0, 0, 0, 0.2);
}

.readonly-user-img {
	width: 100px;
	height: 100px;
	border-radius: 50%;
	object-fit: cover;
	border: 3px solid #667eea;
}
.close-btn {
	background: linear-gradient(135deg, #667eea, #5563c1);
	color: white;
	border: none;
	padding: 10px 20px;
	border-radius: 10px;
	font-weight: bold;
	cursor: pointer;
	transition: transform 0.1s ease-in-out, background 0.3s;
}
.close-btn:hover {
	background: linear-gradient(135deg, #5563c1, #4452a0);
	transform: scale(1.05);
}

</style>

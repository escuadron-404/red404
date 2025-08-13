const postData = {
	id: "123013",
	user: "pixshanghai",
	location: "lechería",
	avatar: "/testUserImage.jpg",
	postMedia: "/testPostImage.jpg",
	description: "going back to the phase where cows are in my head.",
	createdAt: "10h",
	url: "https://red404.app/post/123013",
	likes: 21,
	commentsData: {
		commentCount: 21,
		comments: [
			{
				id: 1,
				author: "leonardo",
				authorAvatar: "/testUserImage.jpg",
				text: "how beautiful",
				isReply: false,
				replies: 1,
				likes: 0,
			},
			{
				id: 2,
				author: "mark",
				authorAvatar: "/testUserImage.jpg",
				text: "󰱱",
				isReply: true,
				replies: null,
				likes: 2,
			},
			{
				id: 3,
				author: "tobias",
				authorAvatar: "/testUserImage.jpg",
				text: "i need to buy some meat",
				isReply: false,
				replies: null,
				likes: 3,
			},
		],
	},
};
export default postData;

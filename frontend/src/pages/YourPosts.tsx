import { useEffect } from "react";

import useAuthStore from "../stores/auth.store";
import usePostStore from "../stores/post.store";

import PostContainer from "../components/Post/PostContainer";

import { H3 } from "../components/Typography";
import PostCreateModal from "../components/Post/PostCreateModal";

const YourPosts = () => {
    const accessToken = useAuthStore((state) => state.accessToken);

    const posts = usePostStore((state) => state.posts);
    const fetchPosts = usePostStore((state) => state.fetchPosts);

    useEffect(() => {
        if (accessToken !== "") {
            fetchPosts();
        }
    }, [accessToken]);

    return (
        <div className="flex flex-col justify-center items-center gap-4 w-4/5 lg:w-3/5 pb-6">
            <PostCreateModal />

            {posts.length > 0 ? (
                posts.map((post) => <PostContainer key={post.id} post={post} />)
            ) : (
                <H3>There's nothing here lol</H3>
            )}
        </div>
    );
};

export default YourPosts;

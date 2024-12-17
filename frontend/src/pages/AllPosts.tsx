import { useEffect } from "react";

import usePostStore from "../stores/post.store";

import PostContainer from "../components/Post/PostContainer";
import PostCreateModal from "../components/Post/PostCreateModal";

import { H3 } from "../components/Typography";

const AllPosts = () => {
    const posts = usePostStore((state) => state.posts);
    const fetchPosts = usePostStore((state) => state.fetchPosts);
    const isLoggedIn = localStorage.getItem("isLoggedIn");

    useEffect(() => {
        fetchPosts();
    }, []);

    return (
        <div className="flex flex-col justify-center items-center gap-4 w-4/5 lg:w-3/5 pb-6">
            {isLoggedIn === "true" ? <PostCreateModal /> : <></>}

            {posts.length > 0 ? (
                posts.map((post) => <PostContainer key={post.id} post={post} />)
            ) : (
                <H3>There's nothing here lol</H3>
            )}
        </div>
    );
};

export default AllPosts;

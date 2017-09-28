import React from 'react';
import {Card, CardHeader, CardTitle, CardText} from 'material-ui/Card';
import axios from 'axios'

export default class ShowCard extends React.Component {
    constructor(props) {
        super(props);
        this.state = { post: {} }
    }

    componentDidMount() {
        this.getPost();
    }

    getPost() {
        axios.get(`http://localhost:4000/posts/${this.props.postId}`)
            .then(res => {
                console.log(res.data.data)
                this.setState({ post: res.data.data });
            });
    }

    render() {
        return (
            <Card>
                <CardHeader
                    title="Bin Joy"
                    subtitle="A Coder"
                    avatar="https://avatars2.githubusercontent.com/u/1658618?v=4&s=460"
                />
                <CardTitle title={this.state.post.title} subtitle="2017/10/8" />
                <CardText>
                    {this.state.post.content}
                </CardText>
            </Card>
        )
    }
}

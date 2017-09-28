import React from 'react';
import {Card, CardActions, CardTitle, CardText} from 'material-ui/Card';
import FlatButton from 'material-ui/FlatButton';
import axios from 'axios'

export default class IndexCard extends React.Component {
    constructor(props) {
        super(props);
        this.state = { posts: [] }
    }

    componentDidMount() {
        this.getIndex();
    }

    getIndex() {
        axios.get(`http://localhost:4000/`)
            .then(res => {
                console.log(res.data.data)
                this.setState({ posts: res.data.data });
            });
    }

    render() {
        return (
            <div>
                { this.state.posts.map(post =>
                    <Card>
                        <CardTitle title={post.title} subtitle="Bin Joy" />
                        <CardText>
                            {post.content}
                        </CardText>
                        <CardActions>
                            <FlatButton label="Details" href={"http://localhost:3000/#/posts/" + post.id} />
                        </CardActions>
                    </Card>)
                }
            </div>
        )
    }
}

import React from 'react'
import ReactDOM from 'react-dom'
import { HashRouter, Route, Link } from 'react-router-dom'
import MuiThemeProvider from 'material-ui/styles/MuiThemeProvider'
import AppBar from 'material-ui/AppBar';

import IndexCard from './index'
import ShowCard from './show'

const App = (props) => {
    return (
        <HashRouter>
            <MuiThemeProvider>
                <div>
                    <AppBar
                        title="GoOnRails"
                        iconClassNameRight="muidocs-icon-navigation-expand-more"
                    />
                    <div>
                        <Route exact path="/" component={Index}/>
                        <Route path="/posts/:id" component={Show}/>
                    </div>
                </div>
            </MuiThemeProvider>
        </HashRouter>
    );
};

const Index = () => {
    return (
        <IndexCard />
    );
};

const Show = ({ match }) => {
    return (
        <ShowCard postId={match.params.id} />
    );
};

document.addEventListener('DOMContentLoaded', () => {
    ReactDOM.render(
        <App />,
        document.body.appendChild(document.createElement('div')),
    )
})

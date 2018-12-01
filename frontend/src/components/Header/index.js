import React from 'react';
import './Header.css'

const axios = require('axios');

export default class Header extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            profile: {
                name: ""
            }
        }
    }

    componentDidMount() {
        const auth = this.props.auth;
        if (auth) {
            axios.get("/api/userinfo")
                .then(response => {
                    let profile = response.data;
                    if (profile.sub.startsWith("vkontakte")) {
                        profile.formatedName = profile.given_name + ' ' + profile.family_name;
                    } else {//if (profile.sub.startsWith("google-oauth2")) {
                        profile.formatedName = profile.name;
                    }

                    this.setState({
                        profile: profile
                    });
                })
        }
    }

    render() {
        const auth = this.props.auth;

        const loggedIn = <section className="navbar-section">
            <div className="name">{this.state.profile.formatedName}</div>
            <div><a href="/logout" className="btn-link bigfont">Logout</a></div>
        </section>;

        const loggedOut = <section className="navbar-section">
            <a href="/login" className="btn-link bigfont">Login</a>
        </section>;

        return (<div>

                <header className="navbar timernavbar bigfont">
                    <section className="navbar-section title">T-timer</section>
                    {auth ? loggedIn : loggedOut}
                </header>
            </div>
        )
    }
}

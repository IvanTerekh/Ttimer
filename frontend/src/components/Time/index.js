import React from "react";
import format from '../../logic/format';
import './Time.css';

function currentDatetime() {
    let d = new Date();
    let s = d.toISOString();
    return s.replace('T', ' ').slice(0, -2);
}

export default class Time extends React.Component {

    constructor(props) {
        super(props);

        this.state = {
            running: false,
            timeState: "usual",
            centis: 0
        };

        this.keydown = this.keydown.bind(this);
        this.keyup = this.keyup.bind(this);
        this.start = this.start.bind(this);
        this.stop = this.stop.bind(this);
        document.onkeydown = (e) => {
            if (e.code === "Space") this.keydown(e)
        };
        document.onkeyup = (e) => {
            if (e.code === "Space") this.keyup(e)
        };
    }

    start() {
        this.setState({
            running: true,
            startTime: Date.now()
        });

        this.timerId = setInterval(
            () => this.tick(), 9);
        document.ontouchstart = this.keydown;
    }

    stop() {
        const res = Math.floor((Date.now() - this.state.startTime) / 10);
        clearInterval(this.timerId);
        this.props.saveResult({
            centis: res,
            scramble: this.props.scramble,
            penalty: false,
            datetime: currentDatetime()
        });
        this.setState({
            centis: res,
            running: false
        });
        document.ontouchstart = () => {
        };
    }

    keydown(e) {
        e.preventDefault();
        if (this.state.running) {
            this.stop();
        } else if (isNaN(this.timeoutId)) {
            this.setState({
                timeState: "notReady"
            });
            this.timeoutId = setTimeout(() => {
                this.setState({
                    timeState: "ready"
                });
            }, 300);
        }
    }

    keyup(e) {
        if (!this.state.running) {
            if (this.state.timeState === "ready") {
                this.start()
            }
            clearTimeout(this.timeoutId);
            this.timeoutId = NaN;

            this.setState({
                timeState: "usual"
            });
        }
    }

    tick() {
        this.setState({
            centis: Math.floor((Date.now() - this.state.startTime) / 10)
        });
    }

    render() {
        return <div onTouchStart={this.keydown} title='Hold "Space" to start.' onTouchEnd={this.keyup}>
            <h1 className={"time " + this.state.timeState}>{format(this.state.centis)}</h1>
        </div>
    }
}

import React from "react";
import Tab from "react-bootstrap/Tab";
import Nav from "react-bootstrap/Nav";

import { RunningTab } from "./tabs/RunningTab";
import { LibrariesTab } from "./tabs/LibrariesTab";
import { HistoryTab } from "./tabs/HistoryTab";
import { SettingsTab } from "./tabs/SettingsTab";
import "./spacers.css";
import "./App.css";

import logo from "./Encodarr-Logo.svg";

function Title() {
	return (<div className="header-flex header-content text-center">
		{/* 100px for Circle; 60px for bare logo; */}
		<img src={logo} alt="" height="60px" title="" />
		<h1>ncodarr</h1>
	</div>);
}

class App extends React.Component {
	handleSelect(eventKey: any) {
		switch (eventKey) {
			case "libraries":
				window.history.replaceState(undefined, "", "/libraries");
				document.title = "Libraries - Encodarr";
				break;
			case "history":
				window.history.replaceState(undefined, "", "/history");
				document.title = "History - Encodarr";
				break;
			case "settings":
				window.history.replaceState(undefined, "", "/settings");
				document.title = "Settings - Encodarr";
				break;
			case "running":
				window.history.replaceState(undefined, "", "/running");
				document.title = "Running - Encodarr";
				break;
			default:
				break;
		}
	}

	render() {
		let activeKey: string = "running";
		switch (window.location.pathname) {
			case "/libraries":
				activeKey = "libraries";
				break;
			case "/history":
				activeKey = "history";
				break;
			case "/settings":
				activeKey = "settings";
				break;
			default:
				break;
		}

		return (<div className="container">
			<Title />
			<Tab.Container id="tab-nav" defaultActiveKey={activeKey} transition={false} onSelect={this.handleSelect}>
				<Nav fill variant="pills">
					<Nav.Item>
						<Nav.Link eventKey="running">Running</Nav.Link>
					</Nav.Item>
					<Nav.Item>
						<Nav.Link eventKey="libraries">Libraries</Nav.Link>
					</Nav.Item>
					<Nav.Item>
						<Nav.Link eventKey="history">History</Nav.Link>
					</Nav.Item>
					<Nav.Item>
						<Nav.Link eventKey="settings">Settings</Nav.Link>
					</Nav.Item>
				</Nav>

				<div className="spacer"></div>

				<Tab.Content>
					<Tab.Pane eventKey="running" mountOnEnter={true} unmountOnExit={true}>
						<RunningTab />
					</Tab.Pane>
					<Tab.Pane eventKey="libraries" mountOnEnter={true} unmountOnExit={true}>
						<LibrariesTab />
					</Tab.Pane>
					<Tab.Pane eventKey="history" mountOnEnter={true} unmountOnExit={true}>
						<HistoryTab />
					</Tab.Pane>
					<Tab.Pane eventKey="settings" mountOnEnter={true} unmountOnExit={true}>
						<SettingsTab />
					</Tab.Pane>
				</Tab.Content>
			</Tab.Container>

			<div className="smol-spacer"></div>
		</div>);
	}
}

export default App;

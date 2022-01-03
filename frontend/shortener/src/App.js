import { useDispatch, useSelector } from "react-redux";
// import Header from "./components/Header";
import UrlForm from "./components/UrlForm";
import UrlHistory from "./components/UrlHistory";
import { addNewURL, createStatus } from "./store/urls/urlSlice";
import { Col, Layout, Menu, Row, Space, Typography, Alert } from "antd";
import {
  HomeOutlined,
  DashboardOutlined,
  TagsOutlined,
} from "@ant-design/icons";
import { logout, user as selectUser } from "./store/user/userSlice";
import { Link } from "react-router-dom";
import { useState } from "react";
import { SUCCEEDED } from "./store/consts";

const { Content, Footer, Sider, Header } = Layout;

function App() {
  const dispatch = useDispatch();
  const user = useSelector(selectUser);
  const newURLStatus = useSelector(createStatus())
  const [currentContent, setCurrentContent] = useState(0);
  const [siderCollapsed, setSiderCollapsed] = useState(false);
  const [alertMessage, setAlertMessage] = useState(null);

  const onUrlShortenClick = (url) => {
    setAlertMessage(null);
    let regexp =
      /^(?:(?:https?|ftp):\/\/)?(?:(?!(?:10|127)(?:\.\d{1,3}){3})(?!(?:169\.254|192\.168)(?:\.\d{1,3}){2})(?!172\.(?:1[6-9]|2\d|3[0-1])(?:\.\d{1,3}){2})(?:[1-9]\d?|1\d\d|2[01]\d|22[0-3])(?:\.(?:1?\d{1,2}|2[0-4]\d|25[0-5])){2}(?:\.(?:[1-9]\d?|1\d\d|2[0-4]\d|25[0-4]))|(?:(?:[a-z\u00a1-\uffff0-9]-*)*[a-z\u00a1-\uffff0-9]+)(?:\.(?:[a-z\u00a1-\uffff0-9]-*)*[a-z\u00a1-\uffff0-9]+)*(?:\.(?:[a-z\u00a1-\uffff]{2,})))(?::\d{2,5})?(?:\/\S*)?$/;
    console.log(url);
    if (!regexp.test(url.url)) {
      setAlertMessage({
        type: "error",
        message: "Invalid URL",
        description: "Given URL is incorrect",
      });
      return;
    }

    dispatch(addNewURL(url));
  };


  const onSignout = () => {
    dispatch(logout());
  };

  const getContent = (contentID) => {
    if (user == null) {
      contentID = 0;
    }
    switch (contentID) {
      case 1:

        return <Typography.Title>Dashboard</Typography.Title>;
      case 2:

        return <Typography.Title>Brands</Typography.Title>;
      case 0:
      default:

        return (
          <Row gutter={[0, 32]} justify="space-around">
            <Col>
              <UrlForm onEnter={onUrlShortenClick} />
            </Col>
            <Col span={12} xs={24} sm={12} md={14} lg={16}>
              <UrlHistory />
            </Col>
          </Row>
        );
    }
  };

  const onMenuClick = ({ key }) => {
    setAlertMessage(null)
    setCurrentContent(key);
  };

  const content = getContent(+currentContent);

  if (newURLStatus === SUCCEEDED) {
    setAlertMessage({
      type: "success",
      message: "URL successfully created"
    })
  }

  return (
    <Layout style={{ minHeight: "100vh" }}>
      {user && (
        <Sider
          style={{
            overflow: "auto",
            height: "100vh",
            position: "fixed",
            left: 0,
          }}
          theme="dark"
          collapsible
          collapsed={siderCollapsed}
          onCollapse={() => setSiderCollapsed(!siderCollapsed)}
        >
          <Space style={{ padding: 8, margin: 16, height: 32 }}>
            <Link to="/profile" component={Typography.Link}>
              {user.username}
            </Link>
          </Space>
          <Menu theme="dark" onClick={onMenuClick}>
            <Menu.Item icon={<HomeOutlined />} key={0}>
              Shortener
            </Menu.Item>
            <Menu.Item icon={<DashboardOutlined />} key={1}>
              Dashboard
            </Menu.Item>
            <Menu.Item icon={<TagsOutlined />} key={2}>
              Brands
            </Menu.Item>
          </Menu>
        </Sider>
      )}
      <Layout style={{ marginLeft: user ? (siderCollapsed ? 50 : 200) : 0 }}>
        <Header>
          <Row gutter={[24, 0]} justify="end">
            <Col span={24} style={{ textAlign: "end" }}>
              {user ? (
                <Typography.Link onClick={onSignout}>Signout</Typography.Link>
              ) : (
                <Link to="/signin" component={Typography.Link}>
                  Signin
                </Link>
              )}
            </Col>
          </Row>
        </Header>
        <Content style={{ paddingTop: "5em", padding: "1em" }}>
          {alertMessage && (
            <Alert
              type={alertMessage.type}
              message={alertMessage.message}
              description={alertMessage.description}
              banner
            />
          )}
          {content}
        </Content>
        <Footer>Copyright IFtech 2021-2022</Footer>
      </Layout>
    </Layout>
  );
}

export default App;

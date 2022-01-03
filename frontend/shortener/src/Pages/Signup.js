import { Form, Input, Button, Checkbox, Avatar, Row, Col } from "antd";
import { UserOutlined } from "@ant-design/icons";
import "antd/dist/antd.css";
import { signup } from "../api/user";

const Signup = () => {
  const onFinish = async (values) => {
    await signup(values);
    console.log("Success", values);
  };
  const onFinishFailed = (errorInfo) => {
    console.log("Failed", errorInfo);
  };
  return (
    <div
      style={{
        display: "flex",
        height: "100vh",
        flexDirection: "column",
        alignItems: "center",
        paddingTop: "15%",
      }}
    >
      <Form
        style={{
          maxWidth: 600,
          width: "50%",
          padding: 16,
          borderRadius: 8,
          boxShadow: "2px 2px 10px 2px rgba(0, 0, 0, 0.05)",
        }}
        name="signin"
        labelCol={{ span: 6 }}
        wrapperCol={{ span: 18 }}
        initialValues={{ remember: true }}
        onFinish={onFinish}
        onFinishFailed={onFinishFailed}
      >
        <Row justify="space-around" align="middle">
          <Col>
            <Avatar
              style={{ marginBottom: 32 }}
              size={54}
              icon={<UserOutlined />}
            />
          </Col>
        </Row>
        <Form.Item
          label="Username"
          name="username"
          rules={[
            {
              required: true,
              message: "Please enter username",
            },
          ]}
        >
          <Input />
        </Form.Item>
        <Form.Item
          label="Password"
          name="password"
          rules={[
            {
              required: true,
              message: "Please enter password",
            },
          ]}
        >
          <Input.Password />
        </Form.Item>
        <Form.Item
          label="Confirm password"
          name="confirm-password"
          rules={[
            {
              required: true,
              message: "Please repeat password",
            },
          ]}
        >
          <Input.Password />
        </Form.Item>

        <Form.Item label="Name" style={{ marginBottom: 0 }}>
          <Form.Item
            name="firstname"
            style={{
              display: "inline-block",
              width: "calc(50% - 8px)",
            }}
            rules={[
              {
                required: true,
                message: "Please enter first name",
              },
            ]}
          >
            <Input placeholder="First name" />
          </Form.Item>
          <Form.Item
            name="lastname"
            style={{
              display: "inline-block",
              width: "calc(50% - 8px)",
              margin: "0 8px",
            }}
            rules={[
              {
                required: true,
                message: "Please enter last name",
              },
            ]}
          >
            <Input placeholder="Last name" />
          </Form.Item>
        </Form.Item>
        <Form.Item
          name="remember"
          valuePropName="checked"
          wrapperCol={{ offset: 6, span: 18 }}
        >
          <Checkbox>Remember Me</Checkbox>
        </Form.Item>

        <Row justify="center" align="middle">
          <Col>
            <Button type="primary" htmlType="submit">
              Submit
            </Button>
          </Col>
        </Row>
      </Form>
    </div>
  );
};

export default Signup;

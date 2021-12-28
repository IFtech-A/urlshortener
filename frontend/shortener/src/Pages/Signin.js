import { Form, Input, Button, Checkbox, Avatar, Row, Col } from "antd";
import { UserOutlined, LockOutlined } from "@ant-design/icons";

import "antd/dist/antd.css";

const Signin = () => {
  const onFinish = (values) => {
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
        name="signin"
        style={{maxWidth: 400, width:'50%', padding:16, borderRadius: 8, boxShadow: '2px 2px 10px 2px rgba(0, 0, 0, 0.05)'}}
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
          <Input prefix={<UserOutlined className="site-form-item-icon" />} />
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
          <Input.Password
            prefix={<LockOutlined className="site-form-item-icon" />}
          />
        </Form.Item>

        <Form.Item
          name="remember"
          valuePropName="checked"
          wrapperCol={{ offset: 6, span: 18 }}
        >
          <Checkbox>Remember Me</Checkbox>
        </Form.Item>

        <Form.Item wrapperCol={{ offset: 6, span: 18 }}>
          <Button type="primary" htmlType="submit">
            Submit
          </Button>
        </Form.Item>
      </Form>
    </div>
  );
};

export default Signin;

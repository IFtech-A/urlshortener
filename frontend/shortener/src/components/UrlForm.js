import { SERVER_HOST } from "../api/consts";
import { Button, Card, Form, Input } from "antd";

const URL = ({ onEnter }) => {
  const onFinish = (values) => {
    console.log(values);
    try {
      onEnter({
        real: values.fullURL,
        short: values.shortURL,
      });
    } catch (e) {
      alert(e);
      console.error(e);
    }
  };
  const onFinishFailed = (errorInfo) => {
    console.log("Failed", errorInfo);
  };

  return (
    <Card
      style={{
        display: "flex",
        justifyContent: "center",
        flexDirection: "column",
        alignContent: "center",
        maxWidth: "500px",
        width: "100%",
      }}
    >
      <Form
        name="url-form"
        wrapperCol={{ span: 24 }}
        onFinish={onFinish}
        onFinishFailed={onFinishFailed}
        initialValues={{ fullURL: "", shortURL: "" }}
      >
        <Form.Item
          name="fullURL"
          rules={[
            {
              require: true,
              message: "Your URL required",
            },
          ]}
          autoFocus
        >
          <Input
            size="large"
            placeholder="Enter your URL"
            type="text"
            width={"100%"}
          />
        </Form.Item>
        <Form.Item name="shortURL">
          <Input size="large" width={"100%"} addonBefore={SERVER_HOST + "/"} />
        </Form.Item>
        <Form.Item>
          <Button type="primary" htmlType="submit">
            Submit
          </Button>
        </Form.Item>
      </Form>
    </Card>
  );
};

export default URL;

import Highlight from "react-highlight";

export type ScheduleValueProps = {
  value: React.ReactNode;
};
const ScheduleValue: React.FC<ScheduleValueProps> = ({ value }) => {
  return <Highlight className="json">{value}</Highlight>;
};

export default ScheduleValue;

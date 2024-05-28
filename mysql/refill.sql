INSERT INTO `record_grouped_rt`
SELECT r.*
FROM record r
         LEFT JOIN record r_next ON r.`rel-id` = r_next.`rel-id` + 1 AND r.file = r_next.file AND r.seq = r_next.seq
         LEFT JOIN record r_prev ON r.`rel-id` = r_prev.`rel-id` - 1 AND r.file = r_prev.file AND r.seq = r_prev.seq
WHERE r.data = 'TCP Retransmission'
  AND (r_next.data = 'TCP Retransmission' OR r_prev.data = 'TCP Retransmission')
ORDER BY r.record_id;
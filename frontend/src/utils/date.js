// src/utils/date.js

/**
 * 格式化日期时间
 * @param {string|Date} date - 日期字符串或Date对象
 * @param {string} format - 格式，默认 'YYYY-MM-DD HH:mm:ss'
 * @returns {string} 格式化后的字符串
 */
export function formatDateTime(date, format = 'YYYY-MM-DD HH:mm:ss') {
  if (!date) return '';
  
  const d = new Date(date);
  if (isNaN(d.getTime())) return '';
  
  const pad = (n) => n.toString().padStart(2, '0');
  
  const replacements = {
    YYYY: d.getFullYear(),
    MM: pad(d.getMonth() + 1),
    DD: pad(d.getDate()),
    HH: pad(d.getHours()),
    mm: pad(d.getMinutes()),
    ss: pad(d.getSeconds()),
  };
  
  return format.replace(/YYYY|MM|DD|HH|mm|ss/g, match => replacements[match]);
}

/**
 * 格式化日期
 * @param {string|Date} date - 日期字符串或Date对象
 * @returns {string} 格式化后的日期字符串
 */
export function formatDate(date) {
  return formatDateTime(date, 'YYYY-MM-DD');
}

/**
 * 格式化时间
 * @param {string|Date} date - 日期字符串或Date对象
 * @returns {string} 格式化后的时间字符串
 */
export function formatTime(date) {
  return formatDateTime(date, 'HH:mm:ss');
}

/**
 * 获取本地化的日期时间字符串
 * @param {string|Date} date - 日期字符串或Date对象
 * @returns {string} 本地化字符串
 */
export function toLocaleDateTime(date) {
  if (!date) return '';
  
  const d = new Date(date);
  if (isNaN(d.getTime())) return '';
  
  return d.toLocaleString('zh-CN');
}

/**
 * 获取本地化的日期字符串
 * @param {string|Date} date - 日期字符串或Date对象
 * @returns {string} 本地化字符串
 */
export function toLocaleDate(date) {
  if (!date) return '';
  
  const d = new Date(date);
  if (isNaN(d.getTime())) return '';
  
  return d.toLocaleDateString('zh-CN');
}
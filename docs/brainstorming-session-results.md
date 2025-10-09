# Brainstorming Session Results

**Session Date:** 2025-10-09
**Facilitator:** Business Analyst Mary
**Participant:** User

---

## Executive Summary

**Topic:** 设计 TMDB MCP 服务 - 确定必要且有意义的接口

**Session Goals:** 广泛探索 TMDB MCP 服务的接口设计可能性，使用 Golang 实现

**Techniques Used:** First Principles Thinking (第一性原理思维) - 30分钟

**Total Ideas Generated:** 10 个原子能力 → 6 个 MCP 工具

### Key Themes Identified:
- **智能化文件管理**：利用 LLM 的泛化能力解决文件重命名难题
- **自然语言交互**：用对话方式探索电影数据库，降低使用门槛
- **场景驱动设计**：从真实痛点出发，而非 API 完整性
- **认知负载优化**：合并相似功能，让 LLM 更容易理解工具

---

## Technique Sessions

### 第一性原理思维 (First Principles Thinking) - 30分钟

**Description:** 从最基本的真理出发，逐层推导，避免类比思维的局限，重新构建问题的解决方案。

#### Ideas Generated:

1. **MCP 的本质价值**
   - LLM 作为"超级胶水"连接现有系统
   - 自然语言成为通用操作接口
   - 从纯对话到实际执行的桥梁

2. **TMDB 核心价值三要素**
   - 稳定的数据来源
   - 元数据更新及时
   - 流媒体软件生态友好

3. **LLM + TMDB 独特场景**
   - 智能文件重命名：混乱文件名 → 标准命名
   - 片荒推荐：自然语言查询高分/热门内容
   - 关联探索：导演/演员作品图谱
   - 智能推荐：基于喜好推荐相似内容

4. **10 个原子能力识别**
   - search_movies, search_tv, search_person, search_multi
   - get_movie_details, get_tv_details, get_person_details
   - discover_movies, discover_tv, get_trending

5. **智能合并策略**
   - 搜索类 4合1 → search(type, query)
   - 详情类 3合1 → get_details(type, id)
   - 发现类保持独立（更清晰的职责）

6. **最终 6 个 MCP 工具设计**
   - search, get_details, discover_movies, discover_tv, get_trending, get_recommendations

#### Insights Discovered:
- **少即是多**：10 个原子能力不等于 10 个工具，合并后更易用
- **LLM 认知优先**：工具数量少比功能完整更重要
- **场景完整性验证**：设计完成后回溯场景，确保覆盖所有需求
- **第一性原理的力量**：从"为什么需要 MCP"到"需要哪些工具"的清晰推导路径

#### Notable Connections:
- TMDB 的"稳定及时"特性 ↔ 自建流媒体服务器的刚需
- LLM 的泛化能力 ↔ 文件名识别的复杂性
- 用户的"片荒"痛点 ↔ discover + trending 的组合使用
- API 原子性 ↔ MCP 工具的合并哲学

---

## Idea Categorization

### Immediate Opportunities
*Ideas ready to implement now*

1. **6 个核心 MCP 工具**
   - Description: search, get_details, discover_movies, discover_tv, get_trending, get_recommendations
   - Why immediate: 基于现有 TMDB API，无需额外开发，直接映射即可实现
   - Resources needed: Golang SDK for TMDB API, MCP 协议实现库

2. **智能文件重命名场景**
   - Description: 使用 search + get_details 解决流媒体文件重命名问题
   - Why immediate: 最高频痛点，ROI 最高
   - Resources needed: 文件系统操作库，LLM 集成测试环境

3. **基础错误处理和速率限制**
   - Description: 实现 401/404/429 错误处理和请求队列
   - Why immediate: TMDB 免费 API 限制 40 requests/10s，必须处理
   - Resources needed: 速率限制库（如 golang.org/x/time/rate）

### Future Innovations
*Ideas requiring development/research*

1. **缓存层设计**
   - Description: 本地缓存热门查询结果，减少 API 调用
   - Development needed: 缓存策略设计、过期机制、持久化选择
   - Timeline estimate: 2-3 周

2. **批量文件处理**
   - Description: 一次性处理整个文件夹的重命名需求
   - Development needed: 批量查询优化、错误恢复机制、进度追踪
   - Timeline estimate: 1-2 周

3. **多语言支持增强**
   - Description: 智能检测用户语言，返回对应语言的元数据
   - Development needed: 语言检测、TMDB language 参数动态设置
   - Timeline estimate: 1 周

4. **高级筛选 DSL**
   - Description: 让 LLM 能将自然语言转换为复杂的 discover 查询
   - Development needed: 研究 LLM 对 discover 参数的理解能力
   - Timeline estimate: 2-4 周

### Moonshots
*Ambitious, transformative concepts*

1. **AI 驱动的个性化推荐引擎**
   - Description: 基于用户历史观看记录，结合 TMDB 推荐 + LLM 理解，生成深度个性化推荐
   - Transformative potential: 超越 TMDB 原生推荐算法，理解用户细微偏好
   - Challenges to overcome: 需要用户数据持久化、隐私保护、LLM 成本控制

2. **跨平台流媒体服务器统一管理**
   - Description: 不仅支持 TMDB，还整合 Jellyfin/Emby/Plex API，实现"一句话管理我的所有影视库"
   - Transformative potential: 成为自建流媒体生态的控制中心
   - Challenges to overcome: 多平台 API 差异、认证管理、实时同步

3. **智能字幕匹配与下载**
   - Description: 结合 TMDB 元数据，自动匹配并下载最佳字幕
   - Transformative potential: 解决字幕查找的最后一公里问题
   - Challenges to overcome: 字幕源整合、语言匹配算法、版权问题

### Insights & Learnings
*Key realizations from the session*

- **第一性原理 > 功能对齐**: 不是"TMDB 有什么 API 我们就暴露什么"，而是"用户需要什么我们就提供什么"
- **合并是艺术**: 4合1 的 search 工具比 4 个独立工具更符合 LLM 的使用习惯
- **场景验证不可少**: 设计完成后必须回到真实场景验证，避免自嗨
- **Golang 的优势**: 静态类型 + 高性能，非常适合实现 MCP 服务的底层逻辑
- **LLM 认知负载**: 工具数量是关键指标，6 个工具是甜蜜点（够用且不过载）

---

## Action Planning

### Top 3 Priority Ideas

#### #1 Priority: 实现 6 个核心 MCP 工具
- **Rationale**: 基础设施优先，没有这些工具，所有场景都无法实现
- **Next steps**:
  1. 搭建 Golang 项目结构，引入 MCP SDK
  2. 实现 TMDB API 客户端封装
  3. 逐个实现 6 个 MCP 工具并编写单元测试
  4. 集成测试（使用真实 API Key）
- **Resources needed**:
  - TMDB API Key
  - Golang MCP SDK（如 github.com/modelcontextprotocol/go-sdk）
  - HTTP 客户端库（net/http 或 resty）
- **Timeline**: 2-3 周

#### #2 Priority: 智能文件重命名场景验证
- **Rationale**: 最高频痛点，验证整个系统的实用价值
- **Next steps**:
  1. 准备测试文件集（包含各种复杂文件名）
  2. 设计 LLM Prompt，引导其使用 search + get_details
  3. 实现文件重命名逻辑（建议先模拟，再真实操作）
  4. 收集反馈并迭代
- **Resources needed**:
  - 测试文件样本
  - Claude/GPT API 集成
  - 文件系统安全操作库
- **Timeline**: 1 周（在工具实现后）

#### #3 Priority: 错误处理和速率限制
- **Rationale**: TMDB 免费 API 有严格限制，必须处理好才能稳定运行
- **Next steps**:
  1. 实现请求队列（限制 40 req/10s）
  2. 实现指数退避重试（处理 429）
  3. 友好的错误信息返回给 LLM
  4. 监控和日志记录
- **Resources needed**:
  - golang.org/x/time/rate
  - 日志库（如 zap）
- **Timeline**: 3-5 天

---

## Reflection & Follow-up

### What Worked Well
- 第一性原理思维让我们避免了"照搬 TMDB API"的陷阱
- 从真实场景出发，确保每个设计决策都有实际价值
- 及时收敛：从 10 个原子能力快速合并到 6 个工具
- 回溯验证：设计完成后验证场景覆盖度

### Areas for Further Exploration
- **参数设计细节**: 每个工具的具体参数列表需要细化（如 discover 的所有筛选条件）
- **响应格式优化**: 如何让返回的 JSON 数据对 LLM 更友好？是否需要简化？
- **图片 URL 处理**: poster_path 需要拼接 base URL，是否在工具层自动完成？
- **多语言策略**: language 参数是每次传入还是全局配置？
- **测试策略**: 如何 mock TMDB API 进行单元测试？

### Recommended Follow-up Techniques
- **Morphological Analysis（形态分析）**: 系统性列出每个工具的所有参数组合，确保覆盖各种使用场景
- **Assumption Reversal（假设逆转）**: 挑战"LLM 一定能正确理解工具"的假设，设计更健壮的错误处理
- **Provocation Technique（挑衅技巧）**: "如果 TMDB API 明天关闭怎么办？"，引导思考备份方案和数据持久化

### Questions That Emerged
- 如何优雅地处理"搜索返回多个结果"的情况？让 LLM 选择还是自动选第一个？
- discover 的参数如此多，LLM 能否有效利用？是否需要提供示例？
- 是否需要提供"获取配置信息"的工具（如支持的语言列表、类型列表）？
- 如何平衡"工具描述的详细程度"与"LLM 的理解负担"？

### Next Session Planning
- **Suggested topics**:
  1. 深入设计每个工具的参数和响应格式
  2. LLM Prompt 工程 - 如何引导 LLM 正确使用这些工具
  3. 实现细节讨论 - Golang 代码架构和最佳实践
- **Recommended timeframe**: 1-2 天后（给予时间消化本次成果）
- **Preparation needed**:
  - 阅读 TMDB API 官方文档的参数详情
  - 调研 Golang MCP SDK 的使用方式
  - 准备一些复杂查询的测试用例

---

*Session facilitated using the BMAD-METHOD™ brainstorming framework*
